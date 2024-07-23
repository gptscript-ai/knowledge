package cohere

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/load"
	cg "github.com/philippgille/chromem-go"
)

const EmbeddingModelProviderCohereName string = "cohere"

type EmbeddingModelProviderCohere struct {
	APIKey string `env:"COHERE_API_KEY" koanf:"apiKey"`
	Model  string `env:"COHERE_EMBEDDING_MODEL" koanf:"model"`
}

func (p *EmbeddingModelProviderCohere) Name() string {
	return EmbeddingModelProviderCohereName
}

func New(configFile string) (*EmbeddingModelProviderCohere, error) {
	p := &EmbeddingModelProviderCohere{}

	err := load.FillConfig(configFile, "COHERE_", &p)
	if err != nil {
		return nil, fmt.Errorf("failed to fill Cohere config")
	}

	if err := p.fillDefaults(); err != nil {
		return nil, fmt.Errorf("failed to fill Cohere defaults: %w", err)
	}

	return p, nil
}

func (p *EmbeddingModelProviderCohere) fillDefaults() error {
	defaultCfg := EmbeddingModelProviderCohere{
		Model: "embed-english-v2.0",
	}

	if err := mergo.Merge(&p, defaultCfg); err != nil {
		return fmt.Errorf("failed to merge Cohere config: %w", err)
	}

	return nil
}

func (p *EmbeddingModelProviderCohere) EmbeddingFunc() (cg.EmbeddingFunc, error) {
	return cg.NewEmbeddingFuncCohere(p.APIKey, cg.EmbeddingModelCohere(p.Model)), nil
}

func (p *EmbeddingModelProviderCohere) Config() any {
	return p
}
