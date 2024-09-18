package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MetadataFilename = ".knowledge.json"

type Metadata struct {
	Metadata map[string]FileMetadata `json:"metadata"` // Map of file paths to metadata
	// TODO (idea): add other fields like description here, so we can hierarchically build a dataset description? Challenge is pruning and merging.
}

type FileMetadata map[string]any

func loadAndMergeMetadata(dirPath string, parentMetadata *Metadata) (*Metadata, error) {
	metadataPath := filepath.Join(dirPath, MetadataFilename)
	dirName := filepath.Base(dirPath)
	if _, err := os.Stat(metadataPath); err == nil { // Metadata file exists
		fileContent, err := os.ReadFile(metadataPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read metadata file %s: %w", metadataPath, err)
		}

		var newMetadata Metadata
		if err := json.Unmarshal(fileContent, &newMetadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata file %s: %w", metadataPath, err)
		}

		// Merge with parent metadata, overriding existing keys
		mergedMetadata := &Metadata{Metadata: map[string]FileMetadata{}}
		for filename, fileMetadata := range parentMetadata.Metadata {
			if !strings.HasPrefix(filename, dirName) {
				// skip entries which are not meant for this (sub-)directory
				continue
			}
			fname := strings.TrimPrefix(strings.TrimPrefix(filename, dirName), string(filepath.Separator))
			mergedMetadata.Metadata[fname] = fileMetadata
		}

		if newMetadata.Metadata != nil {
			for filename, fileMetadata := range newMetadata.Metadata {
				for k, v := range fileMetadata {
					if mergedMetadata.Metadata[filename] == nil {
						mergedMetadata.Metadata[filename] = map[string]any{}
					}
					mergedMetadata.Metadata[filename][k] = v
				}
			}
		}

		return mergedMetadata, nil
	}

	// No metadata file, return parent metadata as is
	return parentMetadata, nil
}
