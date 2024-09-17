package client

import (
	"context"

	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/datastore/textsplitter"
	dstypes "github.com/gptscript-ai/knowledge/pkg/datastore/types"
	"github.com/gptscript-ai/knowledge/pkg/flows"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"github.com/gptscript-ai/knowledge/pkg/server/types"
)

type IngestPathsOpts struct {
	IgnoreExtensions     []string
	Concurrency          int
	Recursive            bool
	TextSplitterOpts     *textsplitter.TextSplitterOpts
	IngestionFlows       []flows.IngestionFlow
	IgnoreFile           string
	IncludeHidden        bool
	NoCreateDataset      bool
	IsDuplicateFuncName  string
	Prune                bool // Prune deleted files
	ErrOnUnsupportedFile bool
}

type Client interface {
	CreateDataset(ctx context.Context, datasetID string) (*index.Dataset, error)
	DeleteDataset(ctx context.Context, datasetID string) error
	GetDataset(ctx context.Context, datasetID string) (*index.Dataset, error)
	ListDatasets(ctx context.Context) ([]types.Dataset, error)
	Ingest(ctx context.Context, datasetID string, name string, data []byte, opts datastore.IngestOpts) ([]string, error)
	IngestPaths(ctx context.Context, datasetID string, opts *IngestPathsOpts, paths ...string) (int, error) // returns number of files ingested
	AskDirectory(ctx context.Context, path string, query string, opts *IngestPathsOpts, ropts *datastore.RetrieveOpts) (*dstypes.RetrievalResponse, error)
	PrunePath(ctx context.Context, datasetID string, path string, keep []string) ([]index.File, error)
	DeleteDocuments(ctx context.Context, datasetID string, documentIDs ...string) error
	Retrieve(ctx context.Context, datasetIDs []string, query string, opts datastore.RetrieveOpts) (*dstypes.RetrievalResponse, error)
	ExportDatasets(ctx context.Context, path string, datasets ...string) error
	ImportDatasets(ctx context.Context, path string, datasets ...string) error
	UpdateDataset(ctx context.Context, dataset index.Dataset, opts *datastore.UpdateDatasetOpts) (*index.Dataset, error)
}
