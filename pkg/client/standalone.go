package client

import (
	"context"
	"fmt"
	dstypes "github.com/gptscript-ai/knowledge/pkg/datastore/types"
	"os"
	"path/filepath"

	"github.com/acorn-io/z"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"github.com/gptscript-ai/knowledge/pkg/server/types"
)

type StandaloneClient struct {
	*datastore.Datastore
}

func NewStandaloneClient(ds *datastore.Datastore) (*StandaloneClient, error) {
	return &StandaloneClient{
		Datastore: ds,
	}, nil
}

func (c *StandaloneClient) CreateDataset(ctx context.Context, datasetID string) (*index.Dataset, error) {
	ds := index.Dataset{
		ID: datasetID,
	}
	err := c.Datastore.NewDataset(ctx, ds)
	if err != nil {
		return &ds, err
	}
	return &ds, nil
}

func (c *StandaloneClient) DeleteDataset(ctx context.Context, datasetID string) error {
	return c.Datastore.DeleteDataset(ctx, datasetID)
}

func (c *StandaloneClient) GetDataset(ctx context.Context, datasetID string) (*index.Dataset, error) {
	return c.Datastore.GetDataset(ctx, datasetID)
}

func (c *StandaloneClient) ListDatasets(ctx context.Context) ([]types.Dataset, error) {
	ds, err := c.Datastore.ListDatasets(ctx)
	if err != nil {
		return nil, err
	}
	r := make([]types.Dataset, len(ds))
	for i, d := range ds {
		r[i] = types.Dataset{
			ID: d.ID,
		}
	}
	return r, nil
}

func (c *StandaloneClient) Ingest(ctx context.Context, datasetID string, data []byte, opts datastore.IngestOpts) ([]string, error) {

	return c.Datastore.Ingest(ctx, datasetID, data, opts)
}

func (c *StandaloneClient) IngestPaths(ctx context.Context, datasetID string, opts *IngestPathsOpts, paths ...string) (int, error) {
	_, err := getOrCreateDataset(ctx, c, datasetID, !opts.NoCreateDataset)
	if err != nil {
		return 0, err
	}

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

		iopts := datastore.IngestOpts{
			Filename: z.Pointer(filepath.Base(path)),
			FileMetadata: &index.FileMetadata{
				Name:         filepath.Base(path),
				AbsolutePath: abspath,
				Size:         finfo.Size(),
				ModifiedAt:   finfo.ModTime(),
			},
			IsDuplicateFunc: datastore.DedupeByFileMetadata,
		}

		if opts != nil {
			iopts.TextSplitterOpts = opts.TextSplitterOpts
			iopts.IngestionFlows = opts.IngestionFlows
		}

		_, err = c.Ingest(ctx, datasetID, file, iopts)
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

func (c *StandaloneClient) Retrieve(ctx context.Context, datasetID string, query string, opts datastore.RetrieveOpts) (*dstypes.RetrievalResponse, error) {
	return c.Datastore.Retrieve(ctx, datasetID, query, opts)
}

func (c *StandaloneClient) AskDirectory(ctx context.Context, path string, query string, opts *IngestPathsOpts, ropts *datastore.RetrieveOpts) (*dstypes.RetrievalResponse, error) {
	return AskDir(ctx, c, path, query, opts, ropts)
}

func (c *StandaloneClient) ExportDatasets(ctx context.Context, path string, datasets ...string) error {
	return c.Datastore.ExportDatasetsToFile(ctx, path, datasets...)
}

func (c *StandaloneClient) ImportDatasets(ctx context.Context, path string, datasets ...string) error {
	return c.Datastore.ImportDatasetsFromFile(ctx, path, datasets...)
}

func (c *StandaloneClient) UpdateDataset(ctx context.Context, dataset index.Dataset, opts *datastore.UpdateDatasetOpts) (*index.Dataset, error) {
	return c.Datastore.UpdateDataset(ctx, dataset, opts)
}
