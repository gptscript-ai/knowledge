package db

import (
	"gorm.io/gorm"
)

type Dataset struct {
	gorm.Model
	Name  string `gorm:"primaryKey"`
	Files []File `gorm:"foreignKey:Dataset;references:Name;constraint:OnDelete:CASCADE;"`
}

type File struct {
	gorm.Model
	FileID    string     `gorm:"primaryKey"`
	Dataset   string     `gorm:"primaryKey"` // Foreign key to Dataset
	Documents []Document `gorm:"gorm:foreignKey:Dataset,FileID;constraint:OnDelete:CASCADE;"`
}

type Document struct {
	gorm.Model
	DocumentID string `gorm:"primaryKey"`
	Dataset    string `gorm:"primaryKey"` // Foreign key to Dataset, part of composite primary key with FileID
	FileID     string `gorm:"primaryKey"` // Foreign key to File, part of composite primary key with Dataset
}
