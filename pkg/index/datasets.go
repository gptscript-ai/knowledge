package index

import (
	"context"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"path/filepath"
)

func GetDataset(db *gorm.DB, id string) (*Dataset, error) {
	var dataset Dataset
	err := db.First(&dataset, "id = ?", id).Error
	return &dataset, err
}

// TODO: Index should become an interface to support this for different databases

func (db *DB) ExportDatasetsToFile(ctx context.Context, path string, ids ...string) error {
	gdb := db.gormDB.WithContext(ctx)

	var datasets []Dataset
	err := gdb.Find(&datasets, "id IN ?", ids).Error
	if err != nil {
		return err
	}

	slog.Debug("Exporting datasets", "ids", ids, "count", len(datasets))

	finfo, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if finfo.IsDir() {
		path = filepath.Join(path, "knowledge-export.db")
	}

	slog.Debug("Exporting datasets to file", "path", path)

	ndb, err := New("sqlite://"+path, true)
	if err != nil {
		return err
	}
	if err := ndb.AutoMigrate(); err != nil {
		return err
	}
	ngdb := ndb.gormDB.WithContext(ctx)

	defer ndb.Close()

	// fill new database with exported datasets
	for _, dataset := range datasets {
		if err := ngdb.Create(&dataset).Error; err != nil {
			return err
		}
	}
	ngdb.Commit()

	return nil
}
