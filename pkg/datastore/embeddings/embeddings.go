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
		name = embeddingsConfig.EmbeddingModelProvider
	}
	embeddingsConfig.EmbeddingModelProvider = name

	switch name {
	case openai.EmbeddingModelProviderOpenAIName:
		return openai.New(embeddingsConfig.EmbeddingModelProviderOpenAI)
	case cohere.EmbeddingModelProviderCohereName:
		return cohere.New(embeddingsConfig.EmbeddingModelProviderCohere)
	case vertex.EmbeddingProviderGoogleVertexAIName:
		return vertex.New(embeddingsConfig.EmbeddingProviderGoogleVertexAI)
	default:
		return nil, fmt.Errorf("unknown embedding model provider: %q", name)
	}
}
