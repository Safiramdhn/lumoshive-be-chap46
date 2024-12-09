package database

import (
	"fmt"
	"golang-chap46/models"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`).Error; err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Define migrations
	allModel := []struct {
		name  string
		model interface{}
	}{
		{"user", models.User{}},
	}

	for _, migration := range allModel {
		var count int64
		err := db.Raw("SELECT COUNT(1) FROM migrations WHERE name = ?", migration.name).Scan(&count).Error
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", migration.name, err)
		}

		if count > 0 {
			log.Printf("Migration '%s' already applied, skipping.", migration.name)
			continue
		}

		// Run migration
		if err := db.AutoMigrate(migration.model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", migration.model, err)
		}

		// Record migration as applied
		if err := db.Exec("INSERT INTO migrations (name) VALUES (?)", migration.name).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.name, err)
		}

		log.Printf("Migration '%s' applied successfully.", migration.name)
	}

	return nil
}
