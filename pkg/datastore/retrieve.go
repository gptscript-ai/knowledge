package datastore

import (
	"context"
	"github.com/gptscript-ai/knowledge/pkg/types/defaults"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"log/slog"
)

func (s *Datastore) Retrieve(ctx context.Context, datasetID string, query string, topk int) ([]vectorstore.Document, error) {
	if topk <= 0 {
		topk = defaults.TopK
	}
	slog.Debug("Retrieving content from dataset", "dataset", datasetID, "query", query)

	docs, err := s.Vectorstore.SimilaritySearch(ctx, query, topk, datasetID)
	if err != nil {
		return nil, err
	}
	slog.Debug("Retrieved documents", "num_documents", len(docs))
	return docs, nil
}
