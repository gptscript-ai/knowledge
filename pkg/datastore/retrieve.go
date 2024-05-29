package datastore

import (
	"context"
	"github.com/gptscript-ai/knowledge/pkg/datastore/defaults"
	"github.com/gptscript-ai/knowledge/pkg/flows"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"log/slog"
)

type RetrieveOpts struct {
	TopK          int
	RetrievalFlow *flows.RetrievalFlow
}

func (s *Datastore) Retrieve(ctx context.Context, datasetID string, query string, opts RetrieveOpts) ([]vectorstore.Document, error) {
	if opts.TopK <= 0 {
		opts.TopK = defaults.TopK
	}
	slog.Debug("Retrieving content from dataset", "dataset", datasetID, "query", query)

	retrievalFlow := opts.RetrievalFlow
	if retrievalFlow == nil {
		retrievalFlow = &flows.RetrievalFlow{}
	}
	retrievalFlow.FillDefaults()

	slog.Debug("Retrieval flow", "flow", *retrievalFlow)

	return retrievalFlow.Run(ctx, s.Vectorstore, query, datasetID)
}
