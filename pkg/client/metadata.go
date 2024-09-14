package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const MetadataFilename = ".knowledge.json"

type Metadata struct {
	Metadata map[string]any `json:"metadata"`
	// TODO (idea): add other fields like description here, so we can hierarchically build a dataset description? Challenge is pruning and merging.
}

func loadAndMergeMetadata(dirPath string, parentMetadata map[string]any) (map[string]any, error) {
	metadataPath := filepath.Join(dirPath, MetadataFilename)
	if _, err := os.Stat(metadataPath); err == nil { // Metadata file exists
		fileContent, err := os.ReadFile(metadataPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read metadata file %s: %w", metadataPath, err)
		}

		var newMetadata map[string]any
		if err := json.Unmarshal(fileContent, &newMetadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata file %s: %w", metadataPath, err)
		}

		// Merge with parent metadata, overriding existing keys
		mergedMetadata := make(map[string]any)
		for k, v := range parentMetadata {
			mergedMetadata[k] = v
		}
		for k, v := range newMetadata {
			mergedMetadata[k] = v
		}

		return mergedMetadata, nil
	}

	// No metadata file, return parent metadata as is
	return parentMetadata, nil
}
