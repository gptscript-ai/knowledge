package index

import "gorm.io/gorm"

func GetDataset(db *gorm.DB, id string) (*Dataset, error) {
	var dataset Dataset
	err := db.First(&dataset, "id = ?", id).Error
	return &dataset, err
}
