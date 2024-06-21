package postprocessors

import (
	"context"
	"github.com/acorn-io/z"
	cohere "github.com/cohere-ai/cohere-go/v2"
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"log/slog"
)

type CohereRerankPostprocessor struct {
	ApiKey string `json:"apiKey" yaml:"apiKey"`
	Model  string
	TopN   int
}

func (c *CohereRerankPostprocessor) Transform(ctx context.Context, query string, docs []vs.Document) ([]vs.Document, error) {
	slog.Debug("Reranking documents", "model", c.Model, "topN", c.TopN, "numDocs", len(docs))
	client := cohereclient.NewClient(cohereclient.WithToken(c.ApiKey))

	docItems := make([]*cohere.RerankRequestDocumentsItem, len(docs))

	for i, doc := range docs {
		docItems[i] = &cohere.RerankRequestDocumentsItem{
			String: doc.Content,
		}
	}

	res, err := client.Rerank(ctx, &cohere.RerankRequest{
		Model:           z.Pointer(c.Model),
		Documents:       docItems,
		Query:           query,
		TopN:            z.Pointer(c.TopN),
		ReturnDocuments: z.Pointer(false),
	})
	if err != nil {
		return nil, err
	}

	rerankedDocs := make([]vs.Document, len(res.Results))

	for i, result := range res.Results {
		rerankedDocs[i] = docs[result.Index]
		slog.Debug("Reranked document", "index", i, "relevanceScore", result.RelevanceScore, "originalIndex", result.Index)

		if len(rerankedDocs[i].Metadata) > 0 {
			rerankedDocs[i].Metadata["rerankRelevanceScore"] = result.RelevanceScore
		} else {
			rerankedDocs[i].Metadata = map[string]interface{}{
				"rerankRelevanceScore": result.RelevanceScore,
			}
		}
	}

	return rerankedDocs, nil
}
