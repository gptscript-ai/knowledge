package tiktoken

import (
	_ "embed"
	"strings"
)

//go:embed resource/gpt2/vocab.bpe
var gpt2Vocab string

//go:embed resource/gpt2/encoder.json
var gpt2Encode string

// NewGPT2 creates a new Codec instance for the GPT-2 tokenization scheme.
// It loads the mergeable ranks from the embedded gpt2Vocab and gpt2Encode resources.
// The function returns a pointer to the Codec or an error if any.
func NewGPT2() (*Codec, error) {
	ranks, err := CovertVocabBPEAndEncoderJSONToMergeableBPERanks(strings.NewReader(gpt2Vocab), strings.NewReader(gpt2Encode))
	if err != nil {
		return nil, err
	}

	return &Codec{
		Name:           "gpt2",
		ExplicitNVocab: 50257,
		PatStr:         `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`,
		MergeableRanks: ranks,
		SpecialTokens: map[string]uint{
			EndOfText: 50256,
		},
	}, nil
}
