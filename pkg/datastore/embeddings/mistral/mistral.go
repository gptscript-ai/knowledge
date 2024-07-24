package mistral

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/load"
	cg "github.com/philippgille/chromem-go"
	"strings"
)

type EmbeddingProviderMistral struct {
	APIKey string `koanf:"apiKey" env:"MISTRAL_API_KEY"`
}

const EmbeddingProviderMistralName = "mistral"

func (p *EmbeddingProviderMistral) Name() string {
	return EmbeddingProviderMistralName
}

func New(c EmbeddingProviderMistral) (*EmbeddingProviderMistral, error) {

	if err := load.FillConfigEnv(strings.ToUpper(EmbeddingProviderMistralName), &c); err != nil {
		return nil, fmt.Errorf("failed to fill Mistral config from environment: %w", err)
	}

	if err := c.fillDefaults(); err != nil {
		return nil, fmt.Errorf("failed to fill Mistral defaults: %w", err)
	}

	return &c, nil
}

func (p *EmbeddingProviderMistral) fillDefaults() error {
	defaultCfg := EmbeddingProviderMistral{}

	if err := mergo.Merge(p, defaultCfg); err != nil {
		return fmt.Errorf("failed to merge Mistral config: %w", err)
	}

	return nil
}

func (p *EmbeddingProviderMistral) EmbeddingFunc() (cg.EmbeddingFunc, error) {
	return cg.NewEmbeddingFuncMistral(p.APIKey), nil
}

func (p *EmbeddingProviderMistral) Config() any {
	return p
}
