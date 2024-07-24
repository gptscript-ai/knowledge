package jina

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/load"
	cg "github.com/philippgille/chromem-go"
	"strings"
)

type EmbeddingProviderJina struct {
	APIKey string `koanf:"apiKey" env:"JINA_API_KEY"`
	Model  string `koanf:"model" env:"JINA_MODEL"`
}

const EmbeddingProviderJinaName = "jina"

func (p *EmbeddingProviderJina) Name() string {
	return EmbeddingProviderJinaName
}

func New(c EmbeddingProviderJina) (*EmbeddingProviderJina, error) {

	if err := load.FillConfigEnv(strings.ToUpper(EmbeddingProviderJinaName), &c); err != nil {
		return nil, fmt.Errorf("failed to fill Jina config from environment: %w", err)
	}

	if err := c.fillDefaults(); err != nil {
		return nil, fmt.Errorf("failed to fill Jina defaults: %w", err)
	}

	return &c, nil
}

func (p *EmbeddingProviderJina) fillDefaults() error {
	defaultCfg := EmbeddingProviderJina{
		Model: "jina-embeddings-v2-base-en",
	}

	if err := mergo.Merge(p, defaultCfg); err != nil {
		return fmt.Errorf("failed to merge Jina config: %w", err)
	}

	return nil
}

func (p *EmbeddingProviderJina) EmbeddingFunc() (cg.EmbeddingFunc, error) {
	return cg.NewEmbeddingFuncJina(p.APIKey, cg.EmbeddingModelJina(p.Model)), nil
}

func (p *EmbeddingProviderJina) Config() any {
	return p
}
