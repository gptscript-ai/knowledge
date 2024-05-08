package client

import (
	"context"
	"fmt"
	"github.com/acorn-io/z"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"github.com/gptscript-ai/knowledge/pkg/types"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"os"
	"path/filepath"
)

type StandaloneClient struct {
	*datastore.Datastore
}

func NewStandaloneClient(ds *datastore.Datastore) (*StandaloneClient, error) {
	return &StandaloneClient{
		Datastore: ds,
	}, nil
}

func (c *StandaloneClient) CreateDataset(ctx context.Context, datasetID string) (types.Dataset, error) {
	ds := types.Dataset{
		ID:             datasetID,
		EmbedDimension: nil,
	}
	err := c.Datastore.NewDataset(ctx, ds)
	if err != nil {
		return ds, err
	}
	return ds, nil
}

func (c *StandaloneClient) DeleteDataset(ctx context.Context, datasetID string) error {
	return c.Datastore.DeleteDataset(ctx, datasetID)
}

func (c *StandaloneClient) GetDataset(ctx context.Context, datasetID string) (*index.Dataset, error) {
	return c.Datastore.GetDataset(ctx, datasetID)
}

func (c *StandaloneClient) ListDatasets(ctx context.Context) ([]types.Dataset, error) {
	return c.Datastore.ListDatasets(ctx)
}

func (c *StandaloneClient) Ingest(ctx context.Context, datasetID string, data []byte, opts datastore.IngestOpts) ([]string, error) {
	return c.Datastore.Ingest(ctx, datasetID, data, opts)
}

func (c *StandaloneClient) IngestPaths(ctx context.Context, datasetID string, opts *IngestPathsOpts, paths ...string) (int, error) {
	ingestFile := func(path string) error {
		// Gather metadata
		finfo, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", path, err)
		}

		abspath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
		}

		file, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		_, err = c.Datastore.Ingest(ctx, datasetID, file, datastore.IngestOpts{
			Filename: z.Pointer(filepath.Base(path)),
			FileMetadata: &index.FileMetadata{
				Name:         filepath.Base(path),
				AbsolutePath: abspath,
				Size:         finfo.Size(),
				ModifiedAt:   finfo.ModTime(),
			},
			IsDuplicateFunc: datastore.DedupeByFileMetadata,
		})
		return err
	}

	return ingestPaths(ctx, opts, ingestFile, paths...)
}

func (c *StandaloneClient) DeleteDocuments(ctx context.Context, datasetID string, documentIDs ...string) error {
	for _, id := range documentIDs {
		err := c.Datastore.DeleteDocument(ctx, datasetID, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *StandaloneClient) Retrieve(ctx context.Context, datasetID string, query string, opts RetrieveOpts) ([]vectorstore.Document, error) {
	return c.Datastore.Retrieve(ctx, datasetID, types.Query{Prompt: query, TopK: z.Pointer(opts.TopK)})
}

func (c *StandaloneClient) AskDirectory(ctx context.Context, path string, query string, opts *IngestPathsOpts, ropts *RetrieveOpts) ([]vectorstore.Document, error) {
	return AskDir(ctx, c, path, query, opts, ropts)
}
