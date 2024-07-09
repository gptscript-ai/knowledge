package querymodifiers

import (
	"context"
	"encoding/json"

	"github.com/gptscript-ai/knowledge/pkg/llm"
)

const SpellcheckQueryModifierName = "spellcheck"

type SpellcheckQueryModifier struct {
	Model llm.LLMConfig
}

func (s SpellcheckQueryModifier) Name() string {
	return SpellcheckQueryModifierName
}

var spellcheckPromptTpl = `The following query will be used for a vector similarity search.
Please spellcheck it and correct any mistakes that may interfere with the semantic similarity search.
Query: "{{.query}}"
Reply only with {"result": "<corrected-query>"}.`

type spellcheckResponse struct {
	Result string `json:"result"`
}

func (s SpellcheckQueryModifier) ModifyQuery(query string) ([]string, error) {
	m, err := llm.NewFromConfig(s.Model)
	if err != nil {
		return nil, err
	}
	result, err := m.Prompt(context.Background(), spellcheckPromptTpl, map[string]interface{}{"query": query})
	if err != nil {
		return nil, err
	}
	var resp spellcheckResponse
	err = json.Unmarshal([]byte(result), &resp)
	if err != nil {
		return nil, err
	}
	return []string{resp.Result}, nil
}
