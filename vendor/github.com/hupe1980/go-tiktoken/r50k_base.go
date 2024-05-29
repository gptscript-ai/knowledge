package tiktoken

import (
	_ "embed"
	"strings"
)

//go:embed resource/r50k_base.tiktoken
var r50kBase string

// NewR50kBase creates a new Codec instance for the R50k_base tokenization scheme.
// It loads the mergeable ranks from the embedded r50kBase resource.
// The function returns a pointer to the Codec or an error if any.
func NewR50kBase() (*Codec, error) {
	ranks, err := ConvertToMergeableBPERanks(strings.NewReader(r50kBase))
	if err != nil {
		return nil, err
	}

	return &Codec{
		Name:           "r50k_base",
		ExplicitNVocab: 50257,
		PatStr:         `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`,
		MergeableRanks: ranks,
		SpecialTokens: map[string]uint{
			EndOfText: 50256,
		},
	}, nil
}
