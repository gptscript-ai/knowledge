package datastore

import (
	"context"

	"github.com/gptscript-ai/knowledge/pkg/index"
)

// IsDuplicateFunc is a function that determines whether a document is a duplicate or if it should be ingested.
// The function should return true if the document is a duplicate (and thus should not be ingested) and false otherwise.
type IsDuplicateFunc func(ctx context.Context, d *Datastore, datasetID string, content []byte, opts IngestOpts) (bool, error)

// IsDuplicateFuncs is a map of deduplication functions by name.
var IsDuplicateFuncs = map[string]IsDuplicateFunc{
	"file_metadata": DedupeByFileMetadata,
	"dummy":         DummyDedupe,
	"none":          DummyDedupe,
	"ignore":        DummyDedupe,
}

// DedupeByFileMetadata is a deduplication function that checks if the document is a duplicate based on the file metadata.
func DedupeByFileMetadata(ctx context.Context, d *Datastore, datasetID string, content []byte, opts IngestOpts) (bool, error) {
	var count int64
	err := d.Index.WithContext(ctx).Model(&index.File{}).
		Where("dataset = ?", datasetID).
		Where("absolute_path = ?", opts.FileMetadata.AbsolutePath).
		Where("size = ?", opts.FileMetadata.Size).
		Where("modified_at = ?", opts.FileMetadata.ModifiedAt).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DummyDedupe is a dummy deduplication function that always returns false (i.e. "No Duplicate").
func DummyDedupe(ctx context.Context, d *Datastore, datasetID string, content []byte, opts IngestOpts) (bool, error) {
	return false, nil
}
