package client

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"os"
	"path/filepath"
)

func ingestPaths(ctx context.Context, opts *IngestPathsOpts, ingestionFunc func(path string) error, paths ...string) error {
	if opts.Concurrency < 1 {
		opts.Concurrency = 10
	}
	sem := semaphore.NewWeighted(int64(opts.Concurrency)) // limit max. concurrency

	g, ctx := errgroup.WithContext(ctx)

	for _, p := range paths {
		path := p
		g.Go(func() error {
			// Wait for a free slot, or exit if the context is done
			if err := sem.Acquire(ctx, 1); err != nil {
				return err
			}
			defer sem.Release(1)

			fileInfo, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("failed to get file info for %s: %w", path, err)
			}

			if fileInfo.IsDir() {
				// Process each file in the directory but skip directories (i.e. don't recurse)
				err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if d.IsDir() {
						return nil
					}
					return ingestionFunc(path)
				})
				if err != nil {
					return err
				}
			} else {
				// Process single file
				err := ingestionFunc(path)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	// Wait for all goroutines in the group to finish
	return g.Wait()
}
