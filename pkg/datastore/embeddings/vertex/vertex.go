package vertex

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/load"
	cg "github.com/philippgille/chromem-go"
)

type EmbeddingProviderGoogleVertexAI struct {
	APIKey      string `koanf:"apiKey" env:"GOOGLE_VERTEX_AI_API_KEY"`
	APIEndpoint string `koanf:"apiEndpoint" env:"GOOGLE_VERTEX_AI_API_ENDPOINT"`
	Project     string `koanf:"project" env:"GOOGLE_VERTEX_AI_PROJECT"`
	Model       string `koanf:"model" env:"GOOGLE_VERTEX_AI_MODEL"`
}

const EmbeddingProviderGoogleVertexAIName = "google_vertex_ai"

func (p *EmbeddingProviderGoogleVertexAI) Name() string {
	return EmbeddingProviderGoogleVertexAIName
}

func New(c EmbeddingProviderGoogleVertexAI) (*EmbeddingProviderGoogleVertexAI, error) {

	if err := load.FillConfigEnv("GOOGLE_VERTEX_AI", &c); err != nil {
		return nil, fmt.Errorf("failed to fill Cohere config from environment: %w", err)
	}

	if err := c.fillDefaults(); err != nil {
		return nil, fmt.Errorf("failed to fill GoogleVertexAI defaults: %w", err)
	}

	return &c, nil
}

func (p *EmbeddingProviderGoogleVertexAI) fillDefaults() error {
	defaultCfg := EmbeddingProviderGoogleVertexAI{
		APIEndpoint: "",
		Project:     "",
		Model:       "text-embedding-004",
	}

	if err := mergo.Merge(p, defaultCfg); err != nil {
		return fmt.Errorf("failed to merge GoogleVertexAI config: %w", err)
	}

	return nil
}

func (p *EmbeddingProviderGoogleVertexAI) EmbeddingFunc() (cg.EmbeddingFunc, error) {
	return cg.NewEmbeddingFuncGoogle(p.APIKey, p.Project, cg.EmbeddingModelGoogle(p.Model), cg.WithGoogleAPIEndpoint(p.APIEndpoint)), nil
}

func (p *EmbeddingProviderGoogleVertexAI) Config() any {
	return p
}
