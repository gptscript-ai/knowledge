package index

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	gormDB      *gorm.DB
	sqlDB       *sql.DB
	autoMigrate bool
}

func New(dsn string, autoMigrate bool) (*DB, error) {
	var (
		gdb   gorm.Dialector
		conns = 1
	)
	if strings.HasPrefix(dsn, "sqlite://") {
		gdb = sqlite.Open(strings.TrimPrefix(dsn, "sqlite://"))
	} else {
		return nil, fmt.Errorf("unsupported database dialect", dsn)
	}
	db, err := gorm.Open(gdb, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.New(log.Default(), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			Colorful:      true,
			LogLevel:      logger.Silent,
		}),
	})
	if err != nil {
		return nil, err
	}

	// Enable foreign key constraint to make sure that deletes cascade
	tx := db.Exec("PRAGMA foreign_keys = ON")
	if tx.Error != nil {
		return nil, tx.Error
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(3 * time.Minute)
	sqlDB.SetMaxIdleConns(conns)
	sqlDB.SetMaxOpenConns(conns)

	return &DB{
		gormDB:      db,
		sqlDB:       sqlDB,
		autoMigrate: autoMigrate,
	}, nil
}

func (db *DB) AutoMigrate() error {
	if !db.autoMigrate {
		return nil
	}

	return db.gormDB.AutoMigrate(
		&Dataset{},
		&File{},
		&Document{},
	)
}

func (db *DB) Check(w http.ResponseWriter, _ *http.Request) {
	if err := db.sqlDB.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

func (db *DB) Close() error {
	return db.sqlDB.Close()
}

func (db *DB) WithContext(ctx context.Context) *gorm.DB {
	return db.gormDB.WithContext(ctx)
}
