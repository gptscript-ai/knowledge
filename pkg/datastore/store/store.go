package store

import (
	"context"
	"github.com/gptscript-ai/knowledge/pkg/index"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type Store interface {
	ListDatasets(ctx context.Context) ([]index.Dataset, error)
	GetDataset(ctx context.Context, datasetID string) (*index.Dataset, error)
	SimilaritySearch(ctx context.Context, query string, numDocuments int, collection string, keywords ...string) ([]vs.Document, error)
}
