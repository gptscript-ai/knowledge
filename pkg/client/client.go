package client

import (
	"context"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/datastore/textsplitter"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"github.com/gptscript-ai/knowledge/pkg/server/types"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type IngestPathsOpts struct {
	IgnoreExtensions []string
	Concurrency      int
	Recursive        bool
	TextSplitterOpts *textsplitter.TextSplitterOpts
}

type RetrieveOpts struct {
	TopK int
}

type Client interface {
	CreateDataset(ctx context.Context, datasetID string) (types.Dataset, error)
	DeleteDataset(ctx context.Context, datasetID string) error
	GetDataset(ctx context.Context, datasetID string) (*index.Dataset, error)
	ListDatasets(ctx context.Context) ([]types.Dataset, error)
	Ingest(ctx context.Context, datasetID string, data []byte, opts datastore.IngestOpts) ([]string, error)
	IngestPaths(ctx context.Context, datasetID string, opts *IngestPathsOpts, paths ...string) (int, error) // returns number of files ingested
	AskDirectory(ctx context.Context, path string, query string, opts *IngestPathsOpts, ropts *RetrieveOpts) ([]vectorstore.Document, error)
	DeleteDocuments(ctx context.Context, datasetID string, documentIDs ...string) error
	Retrieve(ctx context.Context, datasetID string, query string, opts RetrieveOpts) ([]vectorstore.Document, error)
}
