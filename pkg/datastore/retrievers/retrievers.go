package retrievers

import (
	"context"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/store"
	"log/slog"

	"github.com/gptscript-ai/knowledge/pkg/datastore/defaults"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type Retriever interface {
	Retrieve(ctx context.Context, store store.Store, query string, datasetID string) ([]vs.Document, error)
	Name() string
}

func GetRetriever(name string) (Retriever, error) {
	switch name {
	case BasicRetrieverName, "default":
		return &BasicRetriever{TopK: defaults.TopK}, nil
	case SubqueryRetrieverName:
		return &SubqueryRetriever{Limit: 3, TopK: 3}, nil
	case RoutingRetrieverName:
		return &RoutingRetriever{TopK: defaults.TopK}, nil
	default:
		return nil, fmt.Errorf("unknown retriever %q", name)
	}
}

func GetDefaultRetriever() Retriever {
	return &BasicRetriever{TopK: defaults.TopK}
}

const BasicRetrieverName = "basic"

type BasicRetriever struct {
	TopK int
}

func (r *BasicRetriever) Name() string {
	return BasicRetrieverName
}

func (r *BasicRetriever) Retrieve(ctx context.Context, store store.Store, query string, datasetID string) ([]vs.Document, error) {
	log := slog.With("retriever", r.Name())
	if r.TopK <= 0 {
		log.Debug("[BasicRetriever] TopK not set, using default", "default", defaults.TopK)
		r.TopK = defaults.TopK
	}
	return store.SimilaritySearch(ctx, query, r.TopK, datasetID)
}
