package embeddings

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/cohere"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/openai"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/types"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/vertex"
)

func GetEmbeddingsModelProvider(name string, embeddingsConfig config.EmbeddingsConfig) (types.EmbeddingModelProvider, error) {

	if name == "" {
		name = embeddingsConfig.Provider
	}
	embeddingsConfig.Provider = name

	switch name {
	case openai.EmbeddingModelProviderOpenAIName:
		return openai.New(openai.EmbeddingModelProviderOpenAI{OpenAIConfig: embeddingsConfig.OpenAI})
	case cohere.EmbeddingModelProviderCohereName:
		return cohere.New(embeddingsConfig.Cohere)
	case vertex.EmbeddingProviderVertexName:
		return vertex.New(embeddingsConfig.Vertex)
	default:
		return nil, fmt.Errorf("unknown embedding model provider: %q", name)
	}
}
