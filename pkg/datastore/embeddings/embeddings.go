package embeddings

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/openai"
	cg "github.com/philippgille/chromem-go"
)

type EmbeddingModelProvider interface {
	Name() string
	EmbeddingFunc() (cg.EmbeddingFunc, error)
	Config() any
}

func GetEmbeddingsModelProvider(name string, configFile string) (EmbeddingModelProvider, error) {
	switch name {
	case openai.EmbeddingModelProviderOpenAIName:
		return openai.New(configFile)
	default:
		return nil, fmt.Errorf("unknown embedding model provider: %s", name)
	}
}

func NewEmbeddingsFunc(providerName, configFile string) (cg.EmbeddingFunc, error) {
	provider, err := GetEmbeddingsModelProvider(providerName, configFile)
	if err != nil {
		return nil, err
	}

	return provider.EmbeddingFunc()
}
