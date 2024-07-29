package client

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	remotes "github.com/gptscript-ai/knowledge/pkg/datastore/documentloader/remote"
	dstypes "github.com/gptscript-ai/knowledge/pkg/datastore/types"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func isIgnored(ignore gitignore.Matcher, path string) bool {
	return ignore.Match(strings.Split(path, string(filepath.Separator)), false)
}

func checkIgnored(path string, ignoreExtensions []string) bool {
	ext := filepath.Ext(path)
	slog.Debug("checking path", "path", path, "ext", ext, "ignore", ignoreExtensions)
	return slices.Contains(ignoreExtensions, ext)
}

func readIgnoreFile(path string) ([]gitignore.Pattern, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to checkout ignore file %q: %w", path, err)
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("ignore file %q is a directory", path)
	}

	var ps []gitignore.Pattern
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open ignore file %q: %w", path, err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		if !strings.HasPrefix(s, "#") && len(strings.TrimSpace(s)) > 0 {
			ps = append(ps, gitignore.ParsePattern(s, []string{path}))
		}
	}

	return ps, nil
}

func ingestPaths(ctx context.Context, opts *IngestPathsOpts, ingestionFunc func(path string) error, paths ...string) (int, error) {
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
			p := "*." + strings.TrimPrefix(ext, ".")
			ignorePatterns = append(ignorePatterns, gitignore.ParsePattern(p, nil))
		}
	}

	slog.Debug("Ignore patterns", "patterns", ignorePatterns)

	ignore := gitignore.NewMatcher(ignorePatterns)

	if opts.Concurrency < 1 {
		opts.Concurrency = 10
	}
	sem := semaphore.NewWeighted(int64(opts.Concurrency)) // limit max. concurrency

	g, ctx := errgroup.WithContext(ctx)

	for _, p := range paths {
		path := p

		if strings.HasPrefix(path, ".") {
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
				if isIgnored(ignore, subPath) {
					slog.Debug("Ignoring file", "path", subPath, "ignorefile", opts.IgnoreFile, "ignoreExtensions", opts.IgnoreExtensions)
					return nil
				} else {
					slog.Debug("NOT IGNORED file", "path", subPath)
				}

				sp := subPath
				g.Go(func() error {
					if err := sem.Acquire(ctx, 1); err != nil {
						return err
					}
					defer sem.Release(1)

					ingestedFilesCount++
					slog.Debug("Ingesting file", "path", sp)
					return ingestionFunc(sp)
				})
				return nil
			})
			if err != nil {
				return ingestedFilesCount, err
			}
		} else {
			if isIgnored(ignore, path) {
				slog.Debug("Ignoring file", "path", path, "ignorefile", opts.IgnoreFile, "ignoreExtensions", opts.IgnoreExtensions)
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
