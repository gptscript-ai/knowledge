package mixedbread

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/load"
	cg "github.com/philippgille/chromem-go"
	"strings"
)

type EmbeddingProviderMixedbread struct {
	APIKey string `koanf:"apiKey" env:"MIXEDBREAD_API_KEY"`
	Model  string `koanf:"model" env:"MIXEDBREAD_MODEL"`
}

const EmbeddingProviderMixedbreadName = "mixedbread"

func (p *EmbeddingProviderMixedbread) Name() string {
	return EmbeddingProviderMixedbreadName
}

func New(c EmbeddingProviderMixedbread) (*EmbeddingProviderMixedbread, error) {

	if err := load.FillConfigEnv(strings.ToUpper(EmbeddingProviderMixedbreadName), &c); err != nil {
		return nil, fmt.Errorf("failed to fill Mixedbread config from environment: %w", err)
	}

	if err := c.fillDefaults(); err != nil {
		return nil, fmt.Errorf("failed to fill Mixedbread defaults: %w", err)
	}

	return &c, nil
}

func (p *EmbeddingProviderMixedbread) fillDefaults() error {
	defaultCfg := EmbeddingProviderMixedbread{
		Model: "all-MiniLM-L6-v2",
	}

	if err := mergo.Merge(p, defaultCfg); err != nil {
		return fmt.Errorf("failed to merge Mixedbread config: %w", err)
	}

	return nil
}

func (p *EmbeddingProviderMixedbread) EmbeddingFunc() (cg.EmbeddingFunc, error) {
	return cg.NewEmbeddingFuncMixedbread(p.APIKey, cg.EmbeddingModelMixedbread(p.Model)), nil
}

func (p *EmbeddingProviderMixedbread) Config() any {
	return p
}
