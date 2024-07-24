package embeddings

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/cohere"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/jina"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/localai"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/mistral"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/mixedbread"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/ollama"
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
	case jina.EmbeddingProviderJinaName:
		return jina.New(embeddingsConfig.Jina)
	case mistral.EmbeddingProviderMistralName:
		return mistral.New(embeddingsConfig.Mistral)
	case mixedbread.EmbeddingProviderMixedbreadName:
		return mixedbread.New(embeddingsConfig.Mixedbread)
	case localai.EmbeddingProviderLocalAIName:
		return localai.New(embeddingsConfig.LocalAI)
	case ollama.EmbeddingProviderOllamaName:
		return ollama.New(embeddingsConfig.Ollama)
	default:
		return nil, fmt.Errorf("unknown embedding model provider: %q", name)
	}
}
