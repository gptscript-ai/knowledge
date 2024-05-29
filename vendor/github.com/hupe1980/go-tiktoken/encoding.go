package tiktoken

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/dlclark/regexp2"
)

// Encoding represents a text encoding scheme.
type Encoding struct {
	name             string
	specialTokensSet map[string]any
	coreBPE          *coreBPE
}

// NewEncodingByName creates a new Encoding instance based on the given encoding name.
func NewEncodingByName(encoding string) (*Encoding, error) {
	var (
		codec *Codec
		err   error
	)

	switch encoding {
	case CL100kBase:
		codec, err = NewCL100kBase()
	case P50kBase:
		codec, err = NewP50kBase()
	case P50kEdit:
		codec, err = NewP50kEdit()
	case R50kBase:
		codec, err = NewR50kBase()
	case GPT2:
		codec, err = NewGPT2()
	default:
		return nil, fmt.Errorf("unknown encoding: %s", encoding)
	}

	if err != nil {
		return nil, err
	}

	return NewEncoding(codec)
}

// NewEncoding creates a new Encoding instance based on the provided Codec.
func NewEncoding(codec *Codec) (*Encoding, error) {
	coreBPE, err := newCoreBPE(codec.MergeableRanks, codec.SpecialTokens, codec.PatStr)
	if err != nil {
		return nil, err
	}

	specialTokensSet := map[string]any{}
	for k := range codec.SpecialTokens {
		specialTokensSet[k] = true
	}

	return &Encoding{
		name:             codec.Name,
		specialTokensSet: specialTokensSet,
		coreBPE:          coreBPE,
	}, nil
}

// Name returns the name of the Encoding.
func (enc *Encoding) Name() string {
	return enc.name
}

// EncodeOrdinary encodes the given text using the Encoding's core BPE.
func (enc *Encoding) EncodeOrdinary(text string) ([]uint, []string) {
	return enc.coreBPE.EncodeOrdinary(text)
}

var AllSpecial = []string{"all"}

// Encode encodes the given text with the specified allowed and disallowed special tokens.
func (enc *Encoding) Encode(text string, allowedSpecial, disallowedSpecial []string) ([]uint, []string, error) {
	var allowedSpecialSet map[string]any
	if len(allowedSpecial) == 1 && allowedSpecial[0] == "all" {
		allowedSpecialSet = enc.specialTokensSet
	} else {
		allowedSpecialSet = map[string]any{}
		for _, v := range allowedSpecial {
			allowedSpecialSet[v] = nil
		}
	}

	disallowedSpecialSet := map[string]any{}
	for _, v := range disallowedSpecial {
		disallowedSpecialSet[v] = nil
	}

	if len(disallowedSpecial) == 1 && disallowedSpecial[0] == "all" {
		disallowedSpecialSet = difference(enc.specialTokensSet, allowedSpecialSet)
	}

	if len(disallowedSpecialSet) > 0 {
		specialRegex := specialTokenRegex(disallowedSpecialSet)

		m := findRegex2StringMatch(text, specialRegex)
		if m != "" {
			return nil, nil, fmt.Errorf("text contains disallowed special token %s", m)
		}
	}

	ids, tokens := enc.coreBPE.Encode(text, allowedSpecialSet)

	return ids, tokens, nil
}

// Decode decodes the given tokens using the Encoding's core BPE.
func (enc *Encoding) Decode(tokens []uint) []byte {
	return enc.coreBPE.Decode(tokens)
}

// difference calculates the set difference between setA and setB.
func difference(setA, setB map[string]any) map[string]any {
	result := make(map[string]any)

	for k := range setA {
		if _, ok := setB[k]; !ok {
			result[k] = true
		}
	}

	return result
}

// specialTokenRegex generates a regular expression pattern to match disallowed special tokens.
func specialTokenRegex(disallowedSpecialSet map[string]any) *regexp2.Regexp {
	specialRegexStrs := make([]string, 0, len(disallowedSpecialSet))
	for k := range disallowedSpecialSet {
		specialRegexStrs = append(specialRegexStrs, regexp.QuoteMeta(k))
	}

	sort.Strings(specialRegexStrs)

	specialRegex := regexp2.MustCompile(strings.Join(specialRegexStrs, "|"), regexp2.None)

	return specialRegex
}

// findRegex2StringMatch finds the first match of the given regular expression in the text and returns it as a string.
func findRegex2StringMatch(text string, reg *regexp2.Regexp) string {
	m, _ := reg.FindStringMatch(text)
	if m == nil {
		return ""
	}

	return m.String()
}
