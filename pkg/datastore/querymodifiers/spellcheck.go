package querymodifiers

import (
	"context"
	"encoding/json"
	"github.com/gptscript-ai/knowledge/pkg/llm"
)

type SpellcheckQueryModifier struct {
	LLM llm.LLM
}

var promptTpl = `The following query will be used for a vector similarity search.
Please spellcheck it and correct any mistakes that may interfere with the semantic similarity search.
Query: "{.query}"
Reply only with {"result": "<corrected-query>"}.`

type response struct {
	Result string `json:"result"`
}

func (s SpellcheckQueryModifier) ModifyQuery(query string) (string, error) {
	result, err := s.LLM.Prompt(context.Background(), promptTpl, map[string]interface{}{"query": query})
	if err != nil {
		return "", err
	}
	var resp response
	err = json.Unmarshal([]byte(result), &resp)
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}
