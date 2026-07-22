package database

import (
	"blog-backend/config"
	"blog-backend/internal/model"
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg config.Config) {
	if err := os.MkdirAll(filepath.Dir(cfg.Database.Path), 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := DB.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.PostRevision{},
		&model.Tag{},
		&model.Quote{},
		&model.Setting{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Enable WAL mode for better concurrent read performance
	DB.Exec("PRAGMA journal_mode=WAL")

	// Initialize FTS5 full-text search index
	InitFTS()

	log.Println("Database initialized successfully")
}
