package tiktoken

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
)

//go:embed resource/claude.json
var claude string

type claudeJSON struct {
	ExplicitNVocab int             `json:"explicit_n_vocab"`
	PatStr         string          `json:"pat_str"`
	BPERanks       string          `json:"bpe_ranks"`
	SpecialTokens  map[string]uint `json:"special_tokens"`
}

// NewClaude creates a new Codec instance for the claude tokenization scheme.
// It loads the mergeable ranks from the embedded claude resource.
// The function returns a pointer to the Codec or an error if any.
func NewClaude() (*Codec, error) {
	c := claudeJSON{}
	if err := json.Unmarshal([]byte(claude), &c); err != nil {
		return nil, err
	}

	parts := strings.SplitN(c.BPERanks, " ", 3)

	offset, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	tokens := strings.Split(parts[2], " ")

	mergeableRanks := make(map[string]uint, len(tokens))

	for i, token := range tokens {
		t, bErr := base64.StdEncoding.DecodeString(token)
		if bErr != nil {
			return nil, bErr
		}

		mergeableRanks[string(t)] = uint(i * offset)
	}

	return &Codec{
		Name:           "claude",
		ExplicitNVocab: c.ExplicitNVocab,
		PatStr:         c.PatStr,
		MergeableRanks: mergeableRanks,
		SpecialTokens:  c.SpecialTokens,
	}, nil
}
