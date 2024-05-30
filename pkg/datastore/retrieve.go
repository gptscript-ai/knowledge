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

	slog.Debug("Retrieval flow", "flow", *retrievalFlow)

	return retrievalFlow.Run(ctx, s.Vectorstore, query, datasetID)
}
