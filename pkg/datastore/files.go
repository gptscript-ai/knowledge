package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/gptscript-ai/knowledge/pkg/index"
)

// ErrDBFileNotFound is returned when a file is not found.
var ErrDBFileNotFound = errors.New("file not found in database")

func (s *Datastore) DeleteFile(ctx context.Context, datasetID, fileID string) error {
	// Find file in database with associated documents
	var file index.File
	tx := s.Index.WithContext(ctx).Preload("Documents").Where("id = ? AND dataset = ?", fileID, datasetID).First(&file)
	if tx.Error != nil {
		return ErrDBFileNotFound
	}

	// Remove owned documents from VectorStore and Database
	for _, doc := range file.Documents {
		if err := s.Vectorstore.RemoveDocument(ctx, doc.ID, datasetID); err != nil {
			return fmt.Errorf("failed to remove document from VectorStore: %w", err)
		}

		tx = s.Index.WithContext(ctx).Delete(&doc)
		if tx.Error != nil {
			return fmt.Errorf("failed to delete document from DB: %w", tx.Error)
		}
	}

	// Remove file DB
	tx = s.Index.WithContext(ctx).Delete(&file)
	if tx.Error != nil {
		return fmt.Errorf("failed to delete file from DB: %w", tx.Error)
	}

	return nil
}
