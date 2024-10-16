package vectorstore

import (
	"context"
	"fmt"
	"strings"

	etypes "github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/types"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore/chromem"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore/pgvector"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore/types"
	cg "github.com/philippgille/chromem-go"
)

type VectorStore interface {
	CreateCollection(ctx context.Context, collection string) error
	AddDocuments(ctx context.Context, docs []types.Document, collection string) ([]string, error)                                                                                 // @return documentIDs, error
	SimilaritySearch(ctx context.Context, query string, numDocuments int, collection string, where map[string]string, whereDocument []cg.WhereDocument) ([]types.Document, error) //nolint:lll
	RemoveCollection(ctx context.Context, collection string) error
	RemoveDocument(ctx context.Context, documentID string, collection string, where map[string]string, whereDocument []cg.WhereDocument) error
	GetDocuments(ctx context.Context, collection string, where map[string]string, whereDocument []cg.WhereDocument) ([]types.Document, error)

	ImportCollectionsFromFile(ctx context.Context, path string, collections ...string) error
	ExportCollectionsToFile(ctx context.Context, path string, collections ...string) error
}

func New(ctx context.Context, dsn string, embeddingProvider etypes.EmbeddingModelProvider) (VectorStore, error) {
	embeddingFunc, err := embeddingProvider.EmbeddingFunc()
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding function: %w", err)
	}

	s := strings.Split(dsn, "://")
	dialect := s[0]
	dsn = s[1]

	switch s[0] {
	case "chromem":
		return chromem.New(dsn, embeddingFunc)
	case "pgvector":

		return pgvector.New(ctx, dsn, embeddingFunc)
	default:
		return nil, fmt.Errorf("unsupported dialect: %q", dialect)
	}
}
