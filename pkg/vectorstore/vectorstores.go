package vectorstore

import (
	"context"
)

type VectorStore interface {
	CreateCollection(ctx context.Context, collection string) error
	AddDocuments(ctx context.Context, docs []Document, collection string) ([]string, error)                      // @return documentIDs, error
	SimilaritySearch(ctx context.Context, query string, numDocuments int, collection string) ([]Document, error) //nolint:lll
	RemoveCollection(ctx context.Context, collection string) error
	RemoveDocument(ctx context.Context, documentID string, collection string) error
}
