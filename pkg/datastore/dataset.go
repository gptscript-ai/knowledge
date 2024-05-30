package datastore

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gptscript-ai/knowledge/pkg/datastore/defaults"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"gorm.io/gorm"
)

func (s *Datastore) NewDataset(ctx context.Context, dataset index.Dataset) error {
	// Set defaults
	if dataset.EmbedDimension <= 0 {
		dataset.EmbedDimension = defaults.EmbeddingDimension
	}

	// Create dataset
	tx := s.Index.WithContext(ctx).Create(&dataset)
	if tx.Error != nil {
		return tx.Error
	}

	// Create collection
	err := s.Vectorstore.CreateCollection(ctx, dataset.ID)
	if err != nil {
		return err
	}
	slog.Info("Created dataset", "id", dataset.ID)
	return nil
}

func (s *Datastore) DeleteDataset(ctx context.Context, datasetID string) error {
	// Delete dataset
	slog.Info("Deleting dataset", "id", datasetID)
	tx := s.Index.WithContext(ctx).Delete(&index.Dataset{}, "id = ?", datasetID)
	if tx.Error != nil {
		return tx.Error
	}

	// Delete collection
	err := s.Vectorstore.RemoveCollection(ctx, datasetID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Datastore) GetDataset(ctx context.Context, datasetID string) (*index.Dataset, error) {
	// Get dataset with files and associated documents preloaded
	dataset := &index.Dataset{}
	tx := s.Index.WithContext(ctx).Preload("Files.Documents").First(dataset, "id = ?", datasetID)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get dataset %q from DB: %w", datasetID, tx.Error)
	}

	return dataset, nil
}

func (s *Datastore) ListDatasets(ctx context.Context) ([]index.Dataset, error) {
	tx := s.Index.WithContext(ctx).Find(&[]index.Dataset{})
	if tx.Error != nil {
		return nil, tx.Error
	}

	var datasets []index.Dataset
	if err := tx.Scan(&datasets).Error; err != nil {
		return nil, err
	}

	return datasets, nil
}
