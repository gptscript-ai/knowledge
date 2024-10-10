package client

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	remotes "github.com/gptscript-ai/knowledge/pkg/datastore/documentloader/remote"
	dstypes "github.com/gptscript-ai/knowledge/pkg/datastore/types"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func ingestPaths(ctx context.Context, c Client, opts *IngestPathsOpts, datasetID string, ingestionFunc func(path string, metadata map[string]any) error, paths ...string) (int, error) {
	ingestedFilesCount := 0

	var ignorePatterns []gitignore.Pattern
	var err error
	if opts.IgnoreFile != "" {
		ignorePatterns, err = readIgnoreFile(opts.IgnoreFile)
		if err != nil {
			return ingestedFilesCount, fmt.Errorf("failed to read ignore file %q: %w", opts.IgnoreFile, err)
		}
	}

	if len(opts.IgnoreExtensions) > 0 {
		for _, ext := range opts.IgnoreExtensions {
			if ext != "" {
				p := "*." + strings.TrimPrefix(ext, ".")
				ignorePatterns = append(ignorePatterns, gitignore.ParsePattern(p, nil))
			}
		}
	}

	ignorePatterns = append(ignorePatterns, DefaultIgnorePatterns...)

	ignore := gitignore.NewMatcher(ignorePatterns)

	if opts.Concurrency < 1 {
		opts.Concurrency = 10
	}
	sem := semaphore.NewWeighted(int64(opts.Concurrency)) // limit max. concurrency

	g, ctx := errgroup.WithContext(ctx)

	// Stack to store metadata when entering nested directories
	var metadataStack []Metadata

	for _, p := range paths {
		path := p
		var touchedFilePaths []string

		if strings.HasPrefix(filepath.Base(filepath.Clean(path)), ".") {
			if !opts.IncludeHidden {
				slog.Debug("Ignoring hidden path", "path", path)
				continue
			}
		}

		if remotes.IsRemote(path) {
			// Load remote files
			remotePath, err := remotes.LoadRemote(path)
			if err != nil {
				return ingestedFilesCount, fmt.Errorf("failed to load from remote %q: %w", path, err)
			}
			path = remotePath
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			return ingestedFilesCount, fmt.Errorf("failed to get file info for %s: %w", path, err)
		}

		if fileInfo.IsDir() {
			initialMetadata := &Metadata{Metadata: map[string]FileMetadata{}}
			directoryMetadata, err := loadAndMergeMetadata(path, initialMetadata)
			if err != nil {
				return ingestedFilesCount, err
			}
			metadataStack = append(metadataStack, *directoryMetadata)

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

					// One dir level deeper -> load new metadata
					parentMetadata := metadataStack[len(metadataStack)-1]
					newMetadata, err := loadAndMergeMetadata(subPath, &parentMetadata)
					if err != nil {
						return err
					}
					metadataStack = append(metadataStack, *newMetadata)
					return nil
				}
				if isIgnored(ignore, subPath) {
					slog.Debug("Ignoring file", "path", subPath, "ignorefile", opts.IgnoreFile, "ignoreExtensions", opts.IgnoreExtensions)
					return nil
				}

				// Process the file
				sp := subPath
				absPath, err := filepath.Abs(sp)
				if err != nil {
					return fmt.Errorf("failed to get absolute path for %s: %w", sp, err)
				}
				touchedFilePaths = append(touchedFilePaths, absPath)

				currentMetadata := metadataStack[len(metadataStack)-1]

				g.Go(func() error {
					if err := sem.Acquire(ctx, 1); err != nil {
						return err
					}
					defer sem.Release(1)

					ingestedFilesCount++
					slog.Debug("Ingesting file", "path", absPath, "metadata", currentMetadata)
					return ingestionFunc(sp, currentMetadata.Metadata[absPath]) // FIXME: metadata
				})
				return nil
			})
			if err != nil {
				return ingestedFilesCount, err
			}
			// Directory processed, pop metadata
			metadataStack = metadataStack[:len(metadataStack)-1]
		} else {
			if isIgnored(ignore, path) {
				slog.Debug("Ignoring file", "path", path, "ignorefile", opts.IgnoreFile, "ignoreExtensions", opts.IgnoreExtensions)
				continue
			}
			absPath, err := filepath.Abs(path)
			if err != nil {
				return ingestedFilesCount, fmt.Errorf("failed to get absolute path for %s: %w", path, err)
			}
			touchedFilePaths = append(touchedFilePaths, absPath)

			// Process a file directly
			g.Go(func() error {
				if err := sem.Acquire(ctx, 1); err != nil {
					return err
				}
				defer sem.Release(1)

				ingestedFilesCount++
				var fileMetadata FileMetadata
				if len(metadataStack) > 0 {
					currentMetadata := metadataStack[len(metadataStack)-1]
					fileMetadata = currentMetadata.Metadata[filepath.Base(path)]
				}
				return ingestionFunc(path, fileMetadata)
			})
		}

		// Prune files for this basePath
		if opts.Prune && fileInfo.IsDir() {
			g.Go(func() error {
				pruned, err := c.PrunePath(ctx, datasetID, path, touchedFilePaths)
				if err != nil {
					return fmt.Errorf("failed to prune files: %w", err)
				}
				slog.Info("Pruned files", "count", len(pruned), "basePath", path)
				return nil
			})
		}
	}

	// Wait for all goroutines to finish
	return ingestedFilesCount, g.Wait()
}

func HashPath(path string) string {
	hasher := sha1.New()
	hasher.Write([]byte(path))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func AskDir(ctx context.Context, c Client, path string, query string, opts *IngestPathsOpts, ropts *datastore.RetrieveOpts) (*dstypes.RetrievalResponse, error) {
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

	datasetID := HashPath(abspath)
	slog.Debug("Directory Dataset ID hashed", "path", abspath, "id", datasetID)

	_, err = getOrCreateDataset(ctx, c, datasetID, true)
	if err != nil {
		return nil, err
	}

	// ingest files
	if opts == nil {
		opts = &IngestPathsOpts{
			Prune: true,
		}
	}

	ingested, err := c.IngestPaths(ctx, datasetID, opts, path)
	if err != nil {
		return nil, fmt.Errorf("failed to ingest files: %w", err)
	}
	slog.Debug("Ingested files", "count", ingested, "path", abspath)

	// retrieve documents
	return c.Retrieve(ctx, []string{datasetID}, query, *ropts)
}

func getOrCreateDataset(ctx context.Context, c Client, datasetID string, create bool) (*index.Dataset, error) {
	var ds *index.Dataset
	var err error
	ds, err = c.GetDataset(ctx, datasetID)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		if create {
			ds, err = c.CreateDataset(ctx, datasetID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("dataset %q not found", datasetID)
		}
	}
	return ds, nil
}
