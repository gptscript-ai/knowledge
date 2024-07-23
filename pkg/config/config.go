package config

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/cohere"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/openai"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/vertex"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"path"
)

type Config struct {
	EmbeddingsConfig EmbeddingsConfig `koanf:"embeddings" json:"embeddings,omitempty"`
}

type EmbeddingModelProviderConfig struct {
	EmbeddingModelProvider string `usage:"Embedding model provider" default:"openai" env:"KNOW_EMBEDDING_MODEL_PROVIDER" name:"embedding-model-provider" koanf:"provider"`
}

type EmbeddingsConfig struct {
	EmbeddingModelProviderConfig
	openai.EmbeddingModelProviderOpenAI    `koanf:"openai" json:"openai,omitempty"`
	cohere.EmbeddingModelProviderCohere    `koanf:"cohere" json:"cohere,omitempty"`
	vertex.EmbeddingProviderGoogleVertexAI `koanf:"googleVertexAI" json:"googleVertexAI,omitempty"`
}

type DatabaseConfig struct {
	DSN         string `usage:"Server database connection string (default \"sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db\")" default:"" env:"KNOW_DB_DSN"`
	AutoMigrate string `usage:"Auto migrate database" default:"true" env:"KNOW_DB_AUTO_MIGRATE"`
}

type VectorDBConfig struct {
	VectorDBPath string `usage:"VectorDBPath to the vector database (default \"$XDG_DATA_HOME/gptscript/knowledge/vector.db\")" default:"" env:"KNOW_VECTOR_DB_PATH"`
}

func LoadConfig(configFile string) (*Config, error) {
	cfg := &Config{}

	if configFile == "" {
		return cfg, nil
	}

	k := koanf.New(".")
	var pa koanf.Parser
	switch path.Ext(configFile) {
	case ".json":
		pa = json.Parser()
	case ".yaml", ".yml":
		pa = yaml.Parser()
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", path.Ext(configFile))
	}

	if err := k.Load(file.Provider(configFile), pa); err != nil {
		return nil, fmt.Errorf("error loading config file %q: %w", configFile, err)
	}

	if err := k.UnmarshalWithConf("", cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file %q: %w", configFile, err)
	}

	return cfg, nil
}
