package querymodifiers

import (
	"context"
	"encoding/json"
	"github.com/gptscript-ai/knowledge/pkg/llm"
)

const EnhanceQueryModifierName = "enhance"

type EnhanceQueryModifier struct {
	Model llm.LLMConfig
}

func (s EnhanceQueryModifier) Name() string {
	return EnhanceQueryModifierName
}

var enhancePromptTpl = `The following query will be used for a vector similarity search.
Please enhance it to improve the semantic similarity search.
Query: "{{.query}}"
Reply only with {"result": "<enhanced-query>"}.`

type enhanceResp struct {
	Result string `json:"result"`
}

func (s EnhanceQueryModifier) ModifyQuery(query string) ([]string, error) {
	m, err := llm.NewFromConfig(s.Model)
	if err != nil {
		return nil, err
	}
	result, err := m.Prompt(context.Background(), enhancePromptTpl, map[string]interface{}{"query": query})
	if err != nil {
		return nil, err
	}
	var resp enhanceResp
	err = json.Unmarshal([]byte(result), &resp)
	if err != nil {
		return nil, err
	}
	return []string{resp.Result}, nil
}
