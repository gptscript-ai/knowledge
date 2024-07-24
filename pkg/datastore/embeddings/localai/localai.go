package localai

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/load"
	cg "github.com/philippgille/chromem-go"
	"strings"
)

type EmbeddingProviderLocalAI struct {
	Model string `koanf:"model" env:"LOCALAI_MODEL"`
}

const EmbeddingProviderLocalAIName = "localai"

func (p *EmbeddingProviderLocalAI) Name() string {
	return EmbeddingProviderLocalAIName
}

func New(c EmbeddingProviderLocalAI) (*EmbeddingProviderLocalAI, error) {

	if err := load.FillConfigEnv(strings.ToUpper(EmbeddingProviderLocalAIName), &c); err != nil {
		return nil, fmt.Errorf("failed to fill LocalAI config from environment: %w", err)
	}

	if err := c.fillDefaults(); err != nil {
		return nil, fmt.Errorf("failed to fill LocalAI defaults: %w", err)
	}

	return &c, nil
}

func (p *EmbeddingProviderLocalAI) fillDefaults() error {
	defaultCfg := EmbeddingProviderLocalAI{
		Model: "bert-cpp-minilm-v6",
	}

	if err := mergo.Merge(p, defaultCfg); err != nil {
		return fmt.Errorf("failed to merge LocalAI config: %w", err)
	}

	return nil
}

func (p *EmbeddingProviderLocalAI) EmbeddingFunc() (cg.EmbeddingFunc, error) {
	return cg.NewEmbeddingFuncLocalAI(p.Model), nil
}

func (p *EmbeddingProviderLocalAI) Config() any {
	return p
}
