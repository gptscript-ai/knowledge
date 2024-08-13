package datastore

import (
	"context"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	"log/slog"

	"github.com/gptscript-ai/knowledge/pkg/datastore/defaults"
	"github.com/gptscript-ai/knowledge/pkg/flows"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type RetrieveOpts struct {
	TopK          int
	Keywords      []string
	RetrievalFlow *flows.RetrievalFlow
}

func (s *Datastore) Retrieve(ctx context.Context, datasetIDs []string, query string, opts RetrieveOpts) (*types.RetrievalResponse, error) {
	slog.Debug("Retrieving content from dataset", "dataset", datasetIDs, "query", query)

	retrievalFlow := opts.RetrievalFlow
	if retrievalFlow == nil {
		retrievalFlow = &flows.RetrievalFlow{}
	}
	topK := defaults.TopK
	if opts.TopK > 0 {
		topK = opts.TopK
	}
	retrievalFlow.FillDefaults(topK)

	return retrievalFlow.Run(ctx, s, query, datasetIDs, &flows.RetrievalFlowOpts{Keywords: opts.Keywords})
}

func (s *Datastore) SimilaritySearch(ctx context.Context, query string, numDocuments int, datasetID string, keywords ...string) ([]vectorstore.Document, error) {
	return s.Vectorstore.SimilaritySearch(ctx, query, numDocuments, datasetID, keywords...)
}
