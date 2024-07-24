package vertex

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/load"
	cg "github.com/philippgille/chromem-go"
	"strings"
)

type EmbeddingProviderVertex struct {
	APIKey      string `koanf:"apiKey" env:"VERTEX_API_KEY"`
	APIEndpoint string `koanf:"apiEndpoint" env:"VERTEX_API_ENDPOINT"`
	Project     string `koanf:"project" env:"VERTEX_PROJECT"`
	Model       string `koanf:"model" env:"VERTEX_MODEL"`
}

const EmbeddingProviderVertexName = "vertex"

func (p *EmbeddingProviderVertex) Name() string {
	return EmbeddingProviderVertexName
}

func New(c EmbeddingProviderVertex) (*EmbeddingProviderVertex, error) {

	if err := load.FillConfigEnv(strings.ToUpper(EmbeddingProviderVertexName), &c); err != nil {
		return nil, fmt.Errorf("failed to fill Vertex config from environment: %w", err)
	}

	if err := c.fillDefaults(); err != nil {
		return nil, fmt.Errorf("failed to fill Vertex defaults: %w", err)
	}

	return &c, nil
}

func (p *EmbeddingProviderVertex) fillDefaults() error {
	defaultCfg := EmbeddingProviderVertex{
		APIEndpoint: "",
		Project:     "",
		Model:       "text-embedding-004",
	}

	if err := mergo.Merge(p, defaultCfg); err != nil {
		return fmt.Errorf("failed to merge Vertex config: %w", err)
	}

	return nil
}

func (p *EmbeddingProviderVertex) EmbeddingFunc() (cg.EmbeddingFunc, error) {
	cfg := cg.NewVertexConfig(p.APIKey, p.Project, cg.EmbeddingModelVertex(p.Model))
	if p.APIEndpoint != "" {
		cfg = cfg.WithAPIEndpoint(p.APIEndpoint)
	}
	return cg.NewEmbeddingFuncVertex(cfg), nil
}

func (p *EmbeddingProviderVertex) Config() any {
	return p
}
