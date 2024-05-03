package client

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"os"
	"path/filepath"
)

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
