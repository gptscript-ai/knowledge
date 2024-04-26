package client

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
)

func ingestPaths(opts *IngestPathsOpts, ingestionFunc func(path string) error, paths ...string) error {

	ingest := func(path string) error {
		if slices.Contains(opts.IgnoreExtensions, filepath.Ext(path)) {
			slog.Info("Ignoring file based on extension ignore list", "path", path, "ignore_extensions", opts.IgnoreExtensions)
			return nil
		}
		return ingestionFunc(path)
	}

	// Iterate over all paths
	for _, path := range paths {

		fileInfo, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %w", path, err)
		}

		if fileInfo.IsDir() {
			// Read directory contents non-recursively
			err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
				if d.IsDir() {
					return nil
				}
				return ingest(path)
			})
			if err != nil {
				return err
			}
		} else {
			// Read file directly
			err := ingest(path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
