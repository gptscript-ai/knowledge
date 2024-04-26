package client

import (
	"fmt"
	"os"
	"path/filepath"
)

func ingestPaths(ingestionFunc func(path string) error, paths ...string) error {

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
				return ingestionFunc(path)
			})
			if err != nil {
				return err
			}
		} else {
			// Read file directly
			err := ingestionFunc(path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
