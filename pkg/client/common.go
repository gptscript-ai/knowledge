package client

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
)

func checkIgnored(path string, ignoreExtensions []string) bool {
	ext := filepath.Ext(path)
	slog.Debug("checking path", "path", path, "ext", ext, "ignore", ignoreExtensions)
	return slices.Contains(ignoreExtensions, ext)
}

func ingestPaths(ctx context.Context, opts *IngestPathsOpts, ingestionFunc func(path string) error, paths ...string) (int, error) {
	ingestedFilesCount := 0

	if opts.Concurrency < 1 {
		opts.Concurrency = 10
	}
	sem := semaphore.NewWeighted(int64(opts.Concurrency)) // limit max. concurrency

	g, ctx := errgroup.WithContext(ctx)

	for _, p := range paths {
		path := p

		fileInfo, err := os.Stat(path)
		if err != nil {
			return ingestedFilesCount, fmt.Errorf("failed to get file info for %s: %w", path, err)
		}

		if fileInfo.IsDir() {
			// Process directory
			err = filepath.WalkDir(path, func(subPath string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					if subPath == path {
						return nil // Always process the top-level directory
					}
					if !opts.Recursive {
						return filepath.SkipDir // Skip subdirectories if not recursive
					}
					return nil
				}
				if checkIgnored(subPath, opts.IgnoreExtensions) {
					slog.Debug("Skipping ingestion of file", "path", subPath, "reason", "extension ignored")
					return nil
				}

				sp := subPath
				g.Go(func() error {
					if err := sem.Acquire(ctx, 1); err != nil {
						return err
					}
					defer sem.Release(1)

					ingestedFilesCount++
					return ingestionFunc(sp)
				})
				return nil
			})
			if err != nil {
				return ingestedFilesCount, err
			}
		} else {
			if checkIgnored(path, opts.IgnoreExtensions) {
				slog.Debug("Skipping ingestion of file", "path", path, "reason", "extension ignored")
				continue
			}
			// Process a file directly
			g.Go(func() error {
				if err := sem.Acquire(ctx, 1); err != nil {
					return err
				}
				defer sem.Release(1)

				ingestedFilesCount++
				return ingestionFunc(path)
			})
		}
	}

	// Wait for all goroutines to finish
	return ingestedFilesCount, g.Wait()
}

func hashPath(path string) string {
	hasher := sha1.New()
	hasher.Write([]byte(path))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func AskDir(ctx context.Context, c Client, path string, query string, opts *IngestPathsOpts, ropts *RetrieveOpts) ([]vectorstore.Document, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path from %q: %w", path, err)
	}

	finfo, err := os.Stat(abspath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path %q does not exist", abspath)
		}
		return nil, fmt.Errorf("failed to get file info for %q: %w", abspath, err)
	}
	if !finfo.IsDir() {
		return nil, fmt.Errorf("path %q is not a directory", abspath)
	}

	datasetID := hashPath(abspath)
	slog.Debug("Directory Dataset ID hashed", "path", abspath, "id", datasetID)

	// check if dataset exists
	dataset, err := c.GetDataset(ctx, datasetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset %q: %w", datasetID, err)
	}
	if dataset == nil {
		// create dataset
		_, err := c.CreateDataset(ctx, datasetID)
		if err != nil {
			return nil, fmt.Errorf("failed to create dataset %q: %w", datasetID, err)
		}
	}

	// ingest files
	if opts == nil {
		opts = &IngestPathsOpts{}
	}
	ingested, err := c.IngestPaths(ctx, datasetID, opts, path)
	if err != nil {
		return nil, fmt.Errorf("failed to ingest files: %w", err)
	}
	slog.Debug("Ingested files", "count", ingested, "path", abspath)

	// retrieve documents
	return c.Retrieve(ctx, datasetID, query, *ropts)
}
