package datastore

import (
	"context"
	"errors"
	"github.com/acorn-io/z"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"github.com/gptscript-ai/knowledge/pkg/types"
	"github.com/gptscript-ai/knowledge/pkg/types/defaults"
	"gorm.io/gorm"
	"log/slog"
)

func (s *Datastore) NewDataset(ctx context.Context, dataset types.Dataset) error {
	// Set defaults
	if dataset.EmbedDimension == nil || *dataset.EmbedDimension <= 0 {
		dataset.EmbedDimension = z.Pointer(defaults.EmbeddingDimension)
	}

	// Create dataset
	slog.Info("Creating dataset", "id", dataset.ID)
	tx := s.Index.WithContext(ctx).Create(&dataset)
	if tx.Error != nil {
		return tx.Error
	}

	// Create collection
	err := s.Vectorstore.CreateCollection(ctx, dataset.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Datastore) DeleteDataset(ctx context.Context, datasetID string) error {
	// Delete dataset
	slog.Info("Deleting dataset", "id", datasetID)
	tx := s.Index.WithContext(ctx).Delete(&types.Dataset{}, "id = ?", datasetID)
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
	var dataset *index.Dataset
	tx := s.Index.WithContext(ctx).Preload("Files.Documents").First(dataset, "id = ?", datasetID)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return dataset, nil
}

func (s *Datastore) ListDatasets(ctx context.Context) ([]types.Dataset, error) {
	tx := s.Index.WithContext(ctx).Find(&[]types.Dataset{})
	if tx.Error != nil {
		return nil, tx.Error
	}

	var datasets []types.Dataset
	if err := tx.Scan(&datasets).Error; err != nil {
		return nil, err
	}

	return datasets, nil
}
