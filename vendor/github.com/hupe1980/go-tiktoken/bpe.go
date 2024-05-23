package tiktoken

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/dlclark/regexp2"
)

type coreBPE struct {
	encoder              map[string]uint
	decoder              map[uint]string
	specialTokensEncoder map[string]uint
	specialTokensDecoder map[uint]string
	tlRegex              *regexp2.Regexp
	tlSpecialRegex       *regexp2.Regexp
	sortedTokenBytes     [][]byte
}

// newCoreBPE creates a new CoreBPE instance.
// It initializes the CoreBPE with the provided encoder, specialTokensEncoder, and pattern.
func newCoreBPE(encoder map[string]uint, specialTokensEncoder map[string]uint, pattern string) (*coreBPE, error) {
	regex, err := regexp2.Compile(pattern, regexp2.None)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex: %s", err)
	}

	specialRegexStrs := make([]string, 0, len(specialTokensEncoder))
	for k := range specialTokensEncoder {
		specialRegexStrs = append(specialRegexStrs, regexp.QuoteMeta(k))
	}

	specialRegex, err := regexp2.Compile(strings.Join(specialRegexStrs, "|"), regexp2.None)
	if err != nil {
		return nil, fmt.Errorf("error compiling special regex: %s", err)
	}

	decoder := make(map[uint]string, len(encoder))
	for k, v := range encoder {
		decoder[v] = k
	}

	if len(encoder) != len(decoder) {
		return nil, errors.New("encoder and decoder map sizes are different")
	}

	specialTokensDecoder := make(map[uint]string, len(specialTokensEncoder))
	for k, v := range specialTokensEncoder {
		specialTokensDecoder[v] = k
	}

	sortedTokenBytes := make([][]byte, 0, len(encoder))
	for k := range encoder {
		sortedTokenBytes = append(sortedTokenBytes, []byte(k))
	}

	sort.Slice(sortedTokenBytes, func(i, j int) bool {
		return bytes.Compare(sortedTokenBytes[i], sortedTokenBytes[j]) < 0
	})

	return &coreBPE{
		encoder:              encoder,
		specialTokensEncoder: specialTokensEncoder,
		decoder:              decoder,
		specialTokensDecoder: specialTokensDecoder,
		tlRegex:              regex,
		tlSpecialRegex:       specialRegex,
		sortedTokenBytes:     sortedTokenBytes,
	}, nil
}

// Encode performs tokenization and encoding of the input text using the Byte Pair Encoding (BPE) algorithm.
// It takes the input text and a set of allowed special tokens as parameters.
// It returns the encoded token IDs and corresponding tokens as slices.
func (bpe *coreBPE) Encode(text string, allowedSpecial map[string]any) ([]uint, []string) {
	specialRegex := bpe.tlSpecialRegex
	regex := bpe.tlRegex

	retIDs := []uint{}
	retTokens := []string{}

	textLength := len(text)

	start := 0

	for {
		var nextSpecial []int

		startFind := start

		for {
			temp := cutText(text, startFind, textLength)
			nextSpecial = findRegex2StringIndex(temp, specialRegex)

			if nextSpecial != nil {
				token := cutText(text, startFind+nextSpecial[0], startFind+nextSpecial[1])
				if _, ok := allowedSpecial[token]; ok {
					break
				}

				startFind += nextSpecial[1]
			} else {
				break
			}
		}

		end := textLength
		if nextSpecial != nil {
			end = start + nextSpecial[0]
		}

		for _, mat := range findRegex2AllStringMatchIndex(cutText(text, start, end), regex) {
			piece := cutText(text, start+mat[0], start+mat[1])
			if id, ok := bpe.encoder[piece]; ok {
				retIDs = append(retIDs, id)
				retTokens = append(retTokens, piece)

				continue
			}

			ids, tokens := bytePairEncode([]byte(piece), bpe.encoder)
			retIDs = append(retIDs, ids...)
			retTokens = append(retTokens, tokens...)
		}

		if nextSpecial != nil {
			temp := cutText(text, start+nextSpecial[0], start+nextSpecial[1])
			id := bpe.specialTokensEncoder[temp]
			retIDs = append(retIDs, id)
			retTokens = append(retTokens, temp)
			start = start + nextSpecial[1]
		} else {
			break
		}
	}

	return retIDs, retTokens
}

// EncodeOrdinary performs tokenization and encoding of the input text using the Byte Pair Encoding (BPE) algorithm,
// treating all tokens as ordinary tokens (not special tokens).
// It takes the input text as a parameter and returns the encoded token IDs and corresponding tokens as slices.
func (bpe *coreBPE) EncodeOrdinary(text string) ([]uint, []string) {
	retIDs := []uint{}
	retTokens := []string{}

	for _, mat := range findRegex2AllStringMatchIndex(text, bpe.tlRegex) {
		piece := cutText(text, mat[0], mat[1])
		if id, ok := bpe.encoder[piece]; ok {
			retIDs = append(retIDs, id)
			retTokens = append(retTokens, piece)

			continue
		}

		ids, tokens := bytePairEncode([]byte(piece), bpe.encoder)
		retIDs = append(retIDs, ids...)
		retTokens = append(retTokens, tokens...)
	}

	return retIDs, retTokens
}

// Decode performs decoding of the input token IDs and reconstructs the original text.
// It takes the token IDs as a parameter and returns the decoded text as a byte slice.
func (bpe *coreBPE) Decode(tokens []uint) []byte {
	ret := make([]byte, 0, len(tokens)*2)

	for _, token := range tokens {
		tokenBytes, ok := bpe.decoder[token]
		if !ok {
			tokenBytes = bpe.specialTokensDecoder[token]
		}

		if len(tokenBytes) > 0 {
			ret = append(ret, tokenBytes...)
		}
	}

	return ret
}

// bytePairMerge performs the byte pair merging process on the given piece using the provided ranks.
// It returns the merged IDs and tokens. The ranks map should contain precomputed ranks for each token.
func bytePairMerge(piece []byte, ranks map[string]uint) ([]uint, []string) {
	// Structure to store the start index and rank of each part
	type part struct {
		start int
		rank  uint
	}

	// Create initial parts with start index and maximum rank
	parts := make([]part, len(piece)+1)
	for i := 0; i < len(parts); i++ {
		parts[i] = part{i, math.MaxUint}
	}

	// Function to get the rank of a given part
	getRank := func(idx, skip int) uint {
		if idx+skip+2 < len(parts) {
			p := piece[parts[idx].start:parts[idx+skip+2].start]
			if rank, ok := ranks[string(p)]; ok {
				return rank
			}
		}

		return math.MaxUint
	}

	for i := 0; i < len(parts)-2; i++ {
		parts[i].rank = getRank(i, 0)
	}

	for len(parts) > 1 {
		minRank, minIdx := uint(math.MaxUint), 0

		for i, p := range parts[:len(parts)-1] {
			if p.rank < minRank {
				minRank = p.rank
				minIdx = i
			}
		}

		// Break if no minimum rank is found
		if minRank == math.MaxUint {
			break
		}

		parts[minIdx].rank = getRank(minIdx, 1)

		if minIdx > 0 {
			parts[minIdx-1].rank = getRank(minIdx-1, 1)
		}

		parts = append(parts[:minIdx+1], parts[minIdx+2:]...)
	}

	ids := make([]uint, len(parts)-1)
	tokens := make([]string, len(parts)-1)

	for i := 0; i < len(ids); i++ {
		token := string(piece[parts[i].start:parts[i+1].start])
		tokens[i] = token
		ids[i] = ranks[token]
	}

	return ids, tokens
}

// bytePairEncode encodes the given piece using byte pair encoding with the provided ranks.
// It returns the encoded IDs and tokens. The ranks map should contain precomputed ranks for each token.
func bytePairEncode(piece []byte, ranks map[string]uint) ([]uint, []string) {
	if len(piece) == 1 {
		v := ranks[string(piece)]
		return []uint{v}, []string{string(piece)}
	}

	return bytePairMerge(piece, ranks)
}

// findRegex2StringIndex finds the index range of the first occurrence of the regular expression pattern in the given text.
// It returns the index range as a slice [start, end] if a match is found, or nil if no match is found.
// The function takes the input text as a string and the regular expression pattern as a compiled *regexp2.Regexp.
func findRegex2StringIndex(text string, reg *regexp2.Regexp) []int {
	m, _ := reg.FindStringMatch(text)
	if m == nil {
		// If no match is found, return nil.
		return nil
	}

	// Extract the index range of the match and return it as [start, end].
	result := make([]int, 2)
	result[0] = m.Index
	result[1] = m.Index + m.Length

	return result
}

// findRegex2AllStringMatchIndex finds all index ranges of the occurrences of the regular expression pattern in the given text.
// It returns a slice of index ranges, where each range is represented as a slice [start, end].
// The function takes the input text as a string and the regular expression pattern as a compiled *regexp2.Regexp.
func findRegex2AllStringMatchIndex(text string, reg *regexp2.Regexp) [][]int {
	var matches [][]int

	m, _ := reg.FindStringMatch(text)
	for m != nil {
		// Extract the index range of the match and append it to the matches slice.
		result := make([]int, 2)
		result[0] = m.Index
		result[1] = m.Index + m.Length
		matches = append(matches, result)
		m, _ = reg.FindNextMatch(m)
	}

	return matches
}

func cutText(text string, start, end int) string {
	if start < 0 {
		start = 0
	}

	if end > len(text) {
		end = len(text)
	}

	return text[start:end]
}
