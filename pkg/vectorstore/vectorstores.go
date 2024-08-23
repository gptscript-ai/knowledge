package vectorstore

import (
	"context"

	"github.com/philippgille/chromem-go"
)

type VectorStore interface {
	CreateCollection(ctx context.Context, collection string) error
	AddDocuments(ctx context.Context, docs []Document, collection string) ([]string, error)                                                                                      // @return documentIDs, error
	SimilaritySearch(ctx context.Context, query string, numDocuments int, collection string, where map[string]string, whereDocument []chromem.WhereDocument) ([]Document, error) //nolint:lll
	RemoveCollection(ctx context.Context, collection string) error
	RemoveDocument(ctx context.Context, documentID string, collection string, where map[string]string, whereDocument []chromem.WhereDocument) error
	GetDocuments(ctx context.Context, collection string, where map[string]string, whereDocument []chromem.WhereDocument) ([]Document, error)

	ImportCollectionsFromFile(ctx context.Context, path string, collections ...string) error
	ExportCollectionsToFile(ctx context.Context, path string, collections ...string) error
}
