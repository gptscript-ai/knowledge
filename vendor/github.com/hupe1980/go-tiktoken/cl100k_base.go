package tiktoken

import (
	_ "embed"
	"strings"
)

//go:embed resource/cl100k_base.tiktoken
var cl100kBase string

// NewCL100kBase creates a new Codec instance for the cl100k_base tokenization scheme.
// It loads the mergeable ranks from the embedded cl100kBase resource.
// The function returns a pointer to the Codec or an error if any.
func NewCL100kBase() (*Codec, error) {
	ranks, err := ConvertToMergeableBPERanks(strings.NewReader(cl100kBase))
	if err != nil {
		return nil, err
	}

	return &Codec{
		Name:           "cl100k_base",
		PatStr:         `(?i:'s|'t|'re|'ve|'m|'ll|'d)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]+|\s+(?!\S)|\s+`,
		MergeableRanks: ranks,
		SpecialTokens: map[string]uint{
			EndOfText:   100257,
			FimPrefix:   100258,
			FimMiddle:   100259,
			FimSuffix:   100260,
			EndOfPrompt: 100276,
		},
	}, nil
}
