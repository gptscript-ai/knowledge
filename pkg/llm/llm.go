package llm

import (
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/schema"
)

func NewOpenAI(cfg config.OpenAIConfig) (schema.ChatModel, error) {
	return chatmodel.NewOpenAI(cfg.APIKey, func(o *chatmodel.OpenAIOptions) {
		o.BaseURL = cfg.APIBase
		o.ModelName = cfg.Model
	})
}
