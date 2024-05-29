package tiktoken

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
)

// Constants for special tokens.
const (
	StartOfText string = "<|startoftext|>"
	EndOfText   string = "<|endoftext|>"
	FimPrefix   string = "<|fim_prefix|>"
	FimMiddle   string = "<|fim_middle|>"
	FimSuffix   string = "<|fim_suffix|>"
	EndOfPrompt string = "<|endofprompt|>"
)

// Codec represents a token encoding codec.
type Codec struct {
	Name           string          `json:"name"`
	ExplicitNVocab int             `json:"explicit_n_vocab"`
	PatStr         string          `json:"pat_str"`
	MergeableRanks map[string]uint `json:"mergeable_ranks"`
	SpecialTokens  map[string]uint `json:"special_tokens"`
}

// CovertVocabBPEAndEncoderJSONToMergeableBPERanks converts the vocabulary BPE and encoder JSON
// to mergeable BPE ranks.
func CovertVocabBPEAndEncoderJSONToMergeableBPERanks(vocabBPE io.Reader, encoderJSON io.Reader) (map[string]uint, error) {
	rankToIntByte := make([]rune, 0)

	for b := 0; b < 256; b++ {
		if strconv.IsPrint(rune(b)) && rune(b) != rune(' ') {
			rankToIntByte = append(rankToIntByte, rune(b))
		}
	}

	dataGymByteToByte := make(map[string]rune)
	for _, b := range rankToIntByte {
		dataGymByteToByte[string(b)] = b
	}

	n := 0

	for b := 0; b < 256; b++ {
		if !containsRune(rankToIntByte, rune(b)) {
			rankToIntByte = append(rankToIntByte, rune(b))
			dataGymByteToByte[string(rune(256+n))] = rune(b)
			n++
		}
	}

	if len(rankToIntByte) != 256 {
		return nil, errors.New("assertion failed: len(rankToIntByte) != 256")
	}

	vocabBPEContents, err := io.ReadAll(vocabBPE)
	if err != nil {
		return nil, err
	}

	vocabBPELines := strings.Split(string(vocabBPEContents), "\n")

	bpeMerges := make([][]string, 0)

	for _, mergeStr := range vocabBPELines[1 : len(vocabBPELines)-1] {
		merge := strings.Split(mergeStr, " ")
		bpeMerges = append(bpeMerges, merge)
	}

	// add the single byte tokens
	bpeRanks := make(map[string]uint)

	for i, b := range rankToIntByte {
		key := string(b)
		bpeRanks[key] = uint(i)
	}

	// add the merged tokens
	decodeDataGym := func(value string) []rune {
		result := []rune{}
		for _, c := range value {
			result = append(result, dataGymByteToByte[string(c)])
		}

		return result
	}

	n = len(bpeRanks)

	for _, merge := range bpeMerges {
		first := decodeDataGym(merge[0])
		second := decodeDataGym(merge[1])
		key := string(append(first, second...))
		bpeRanks[key] = uint(n)
		n++
	}

	encoderJSONContens, err := io.ReadAll(encoderJSON)
	if err != nil {
		return nil, err
	}

	encoderMap := make(map[string]interface{})

	err = json.Unmarshal(encoderJSONContens, &encoderMap)
	if err != nil {
		return nil, err
	}

	encoderLoaded := make(map[string]uint)

	for k, v := range encoderMap {
		result := []rune{}
		for _, r := range k {
			result = append(result, dataGymByteToByte[string(r)])
		}

		encoderLoaded[string(result)] = uint(v.(float64))
	}

	// delete these two special tokens if present, since they're not
	// mergeable bpe tokens
	delete(encoderLoaded, EndOfText)
	delete(encoderLoaded, StartOfText)

	if len(bpeRanks) != len(encoderLoaded) {
		return nil, errors.New("assertion failed: len(bpeRanks) != len(encoderLoaded)")
	}

	for k, v := range bpeRanks {
		if encoderLoaded[k] != v {
			return nil, errors.New("assertion failed: bpeRanks[k] != encoderLoaded[k]")
		}
	}

	return bpeRanks, nil
}

// ConvertToMergeableBPERanks converts the BPE file to mergeable BPE ranks.
func ConvertToMergeableBPERanks(bpe io.Reader) (map[string]uint, error) {
	contents, err := io.ReadAll(bpe)
	if err != nil {
		return nil, err
	}

	if len(contents) == 0 {
		return nil, errors.New("empty bpe file")
	}

	bpeRanks := make(map[string]uint)

	for _, line := range strings.Split(string(contents), "\n") {
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")

		token, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			return nil, err
		}

		rank, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		bpeRanks[string(token)] = uint(rank)
	}

	return bpeRanks, nil
}

// containsRune checks if a rune exists in an array of runes.
func containsRune(arr []rune, b rune) bool {
	for _, v := range arr {
		if v == b {
			return true
		}
	}

	return false
}
