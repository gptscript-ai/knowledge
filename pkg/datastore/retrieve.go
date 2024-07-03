package datastore

import (
	"context"
	"log/slog"

	"github.com/gptscript-ai/knowledge/pkg/datastore/defaults"
	"github.com/gptscript-ai/knowledge/pkg/flows"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type RetrieveOpts struct {
	TopK          int
	RetrievalFlow *flows.RetrievalFlow
}

func (s *Datastore) Retrieve(ctx context.Context, datasetID string, query string, opts RetrieveOpts) ([]vectorstore.Document, error) {
	slog.Debug("Retrieving content from dataset", "dataset", datasetID, "query", query)

	retrievalFlow := opts.RetrievalFlow
	if retrievalFlow == nil {
		retrievalFlow = &flows.RetrievalFlow{}
	}
	topK := defaults.TopK
	if opts.TopK > 0 {
		topK = opts.TopK
	}
	retrievalFlow.FillDefaults(topK)

	return retrievalFlow.Run(ctx, s, query, datasetID)
}

func (s *Datastore) SimilaritySearch(ctx context.Context, query string, numDocuments int, datasetID string) ([]vectorstore.Document, error) {
	return s.Vectorstore.SimilaritySearch(ctx, query, numDocuments, datasetID)
}
