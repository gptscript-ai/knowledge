package tiktoken

import (
	_ "embed"
	"strings"
)

//go:embed resource/p50k_base.tiktoken
var p50kBase string

// NewP50kBase creates a new Codec instance for the P50k_base tokenization scheme.
// It loads the mergeable ranks from the embedded p50kBase resource.
// The function returns a pointer to the Codec or an error if any.
func NewP50kBase() (*Codec, error) {
	ranks, err := ConvertToMergeableBPERanks(strings.NewReader(p50kBase))
	if err != nil {
		return nil, err
	}

	return &Codec{
		Name:           "p50k_base",
		ExplicitNVocab: 50281,
		PatStr:         `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`,
		MergeableRanks: ranks,
		SpecialTokens: map[string]uint{
			EndOfText: 50256,
		},
	}, nil
}
