package types

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/gptscript-ai/knowledge/pkg/config"
	cg "github.com/philippgille/chromem-go"
)

type EmbeddingModelProvider interface {
	Name() string
	EmbeddingFunc() (cg.EmbeddingFunc, error)
	Configure() error
	Config() any
}

func AsEmbeddingModelProviderConfig(emp EmbeddingModelProvider) (config.EmbeddingsProviderConfig, error) {

	var cfg map[string]any

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "koanf",
		Result:  &cfg,
	})

	if err != nil {
		return config.EmbeddingsProviderConfig{}, err
	}

	if err := decoder.Decode(emp.Config()); err != nil {
		return config.EmbeddingsProviderConfig{}, err
	}

	return config.EmbeddingsProviderConfig{
		Type:   emp.Name(),
		Config: cfg,
	}, nil
}
