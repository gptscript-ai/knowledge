package db

// Dataset refers to a VectorDB data space.
// @Description Dataset refers to a VectorDB data space.
type Dataset struct {
	ID             string `gorm:"primaryKey" json:"id"`
	EmbedDimension int    `json:"embed_dim,omitempty"`
	Files          []File `gorm:"foreignKey:Dataset;references:ID;constraint:OnDelete:CASCADE;"`
}

type File struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	Dataset   string     `gorm:"primaryKey" json:"dataset"` // Foreign key to Dataset
	Documents []Document `gorm:"foreignKey:FileID,Dataset;references:ID,Dataset;constraint:OnDelete:CASCADE;"`
}

type Document struct {
	ID      string `gorm:"primaryKey" json:"id"`
	Dataset string `gorm:"primaryKey" json:"dataset"` // Foreign key to Dataset, part of composite primary key with FileID
	FileID  string `gorm:"primaryKey" json:"file_id"` // Foreign key to File, part of composite primary key with Dataset
}
