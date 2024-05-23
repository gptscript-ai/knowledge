package tiktoken

import (
	_ "embed"
	"strings"
)

// NewP50kEdit creates a new Codec instance for the P50k_edit tokenization scheme.
// It loads the mergeable ranks from the embedded p50kBase resource.
// The function returns a pointer to the Codec or an error if any.
func NewP50kEdit() (*Codec, error) {
	ranks, err := ConvertToMergeableBPERanks(strings.NewReader(p50kBase))
	if err != nil {
		return nil, err
	}

	return &Codec{
		Name:           "p50k_edit",
		PatStr:         `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`,
		MergeableRanks: ranks,
		SpecialTokens: map[string]uint{
			EndOfText: 50256,
			FimPrefix: 50281,
			FimMiddle: 50282,
			FimSuffix: 50283,
		},
	}, nil
}
