package db

import (
	"gorm.io/gorm"
)

type Dataset struct {
	gorm.Model
	Name  string      `gorm:"primaryKey"`
	Files []FileIndex `gorm:"constraint:OnDelete:CASCADE;"`
}

type FileIndex struct {
	gorm.Model
	FileID    string          `gorm:"primaryKey"`
	Dataset   string          `gorm:"primaryKey"` // Foreign key to Dataset
	Documents []DocumentIndex `gorm:"constraint:OnDelete:CASCADE;"`
}

type DocumentIndex struct {
	gorm.Model
	DocumentID string `gorm:"primaryKey"`
	Dataset    string `gorm:"primaryKey"` // Foreign key to Dataset, part of composite primary key with FileID
	FileID     string `gorm:"primaryKey"` // Foreign key to FileIndex, part of composite primary key with Dataset
}
