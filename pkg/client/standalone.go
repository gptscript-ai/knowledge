package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader"
	dstypes "github.com/gptscript-ai/knowledge/pkg/datastore/types"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"github.com/gptscript-ai/knowledge/pkg/log"
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

func (c *StandaloneClient) FindFile(ctx context.Context, searchFile index.File) (*index.File, error) {
	return c.Datastore.FindFile(ctx, searchFile)
}

func (c *StandaloneClient) DeleteFile(ctx context.Context, datasetID, fileID string) error {
	return c.Datastore.DeleteFile(ctx, datasetID, fileID)
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

func (c *StandaloneClient) Ingest(ctx context.Context, datasetID string, name string, data []byte, opts datastore.IngestOpts) ([]string, error) {
	ids, err := c.Datastore.Ingest(ctx, datasetID, name, data, opts)
	if err != nil {
		log.FromCtx(ctx).With("status", "failed").With("error", err.Error()).Error("Ingest failed")
	}
	return ids, err
}

func (c *StandaloneClient) IngestPaths(ctx context.Context, datasetID string, opts *IngestPathsOpts, paths ...string) (int, error) {
	_, err := getOrCreateDataset(ctx, c, datasetID, !opts.NoCreateDataset)
	if err != nil {
		return 0, err
	}

	ingestFile := func(path string, extraMetadata map[string]any) error {
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

		filename := filepath.Base(path)

		iopts := datastore.IngestOpts{
			FileMetadata: &index.FileMetadata{
				Name:         filepath.Base(path),
				AbsolutePath: abspath,
				Size:         finfo.Size(),
				ModifiedAt:   finfo.ModTime(),
			},
			IsDuplicateFuncName: opts.IsDuplicateFuncName,
			ExtraMetadata:       extraMetadata,
		}

		if opts != nil {
			iopts.TextSplitterOpts = opts.TextSplitterOpts
			iopts.IngestionFlows = opts.IngestionFlows
		}

		_, err = c.Ingest(log.ToCtx(ctx, log.FromCtx(ctx).With("filepath", path).With("absolute_path", iopts.FileMetadata.AbsolutePath)), datasetID, filename, file, iopts)

		if err != nil && !opts.ErrOnUnsupportedFile && errors.Is(err, &documentloader.UnsupportedFileTypeError{}) {
			err = nil
		}

		return err
	}

	return ingestPaths(ctx, c, opts, datasetID, ingestFile, paths...)
}

func (c *StandaloneClient) PrunePath(ctx context.Context, datasetID string, path string, keep []string) ([]index.File, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}
	return c.Datastore.PruneFiles(ctx, datasetID, abs, keep)
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

func (c *StandaloneClient) Retrieve(ctx context.Context, datasetIDs []string, query string, opts datastore.RetrieveOpts) (*dstypes.RetrievalResponse, error) {
	return c.Datastore.Retrieve(ctx, datasetIDs, query, opts)
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
