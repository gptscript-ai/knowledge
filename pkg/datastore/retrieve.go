package datastore

import (
	"context"
	"github.com/acorn-io/z"
	"github.com/gptscript-ai/knowledge/pkg/types"
	"github.com/gptscript-ai/knowledge/pkg/types/defaults"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"log/slog"
)

func (s *Datastore) Retrieve(ctx context.Context, datasetID string, query types.Query) ([]vectorstore.Document, error) {
	if query.TopK == nil {
		query.TopK = z.Pointer(defaults.TopK)
	}
	slog.Debug("Retrieving content from dataset", "dataset", datasetID, "query", query)

	docs, err := s.Vectorstore.SimilaritySearch(ctx, query.Prompt, *query.TopK, datasetID)
	if err != nil {
		return nil, err
	}
	slog.Debug("Retrieved documents", "num_documents", len(docs))
	return docs, nil
}
