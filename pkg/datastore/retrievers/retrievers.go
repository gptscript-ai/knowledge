package retrievers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gptscript-ai/knowledge/pkg/datastore/defaults"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type Retriever interface {
	Retrieve(ctx context.Context, store vs.VectorStore, query string, datasetID string) ([]vs.Document, error)
}

func GetRetriever(name string) (Retriever, error) {
	switch name {
	case "basic", "default":
		return &BasicRetriever{TopK: defaults.TopK}, nil
	case "subquery":
		return &SubqueryRetriever{Limit: 3, TopK: 3}, nil
	default:
		return nil, fmt.Errorf("unknown retriever %q", name)
	}
}

func GetDefaultRetriever() Retriever {
	return &BasicRetriever{TopK: defaults.TopK}
}

type BasicRetriever struct {
	TopK int
}

func (r *BasicRetriever) Retrieve(ctx context.Context, store vs.VectorStore, query string, datasetID string) ([]vs.Document, error) {
	if r.TopK <= 0 {
		slog.Debug("[BasicRetriever] TopK not set, using default", "default", defaults.TopK)
		r.TopK = defaults.TopK
	}
	return store.SimilaritySearch(ctx, query, r.TopK, datasetID)
}
