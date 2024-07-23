package embeddings

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/cohere"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/openai"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/types"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/vertex"
	cg "github.com/philippgille/chromem-go"
)

func GetEmbeddingsModelProvider(name string, configFile string) (types.EmbeddingModelProvider, error) {
	switch name {
	case openai.EmbeddingModelProviderOpenAIName:
		return openai.New(configFile)
	case cohere.EmbeddingModelProviderCohereName:
		return cohere.New(configFile)
	case vertex.EmbeddingProviderGoogleVertexAIName:
		return vertex.New(configFile)
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
