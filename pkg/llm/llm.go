package llm

import (
	"context"

	"github.com/gptscript-ai/knowledge/pkg/config"
	golcmodel "github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

type LLM struct {
	model schema.Model
}

func NewOpenAI(cfg config.OpenAIConfig) (*LLM, error) {
	m, err := chatmodel.NewOpenAI(cfg.APIKey, func(o *chatmodel.OpenAIOptions) {
		o.BaseURL = cfg.APIBase
		o.ModelName = cfg.Model
	})
	if err != nil {
		return nil, err
	}

	return &LLM{model: m}, nil
}

func (llm *LLM) Prompt(ctx context.Context, promptTpl string, values map[string]any) (string, error) {
	p, err := prompt.NewTemplate(promptTpl).FormatPrompt(values)
	if err != nil {
		return "", err
	}

	res, err := golcmodel.GeneratePrompt(ctx, llm.model, p)
	if err != nil {
		return "", err
	}

	return res.Generations[0].Text, nil
}
