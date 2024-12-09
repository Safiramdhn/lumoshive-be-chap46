package database

import (
	"fmt"
	"golang-chap46/models"
	"log"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		seeds := dataSeeds()
		for _, seed := range seeds {
			var count int64
			name := reflect.TypeOf(seed).String()

			if err := tx.Model(seed).Count(&count).Error; err != nil {
				log.Fatalf("Error checking data for table %s: %v", name, err)
				return err
			}

			if count > 0 {
				log.Printf("Seeding skipped for table %s, data already exists.", name)
				continue
			}

			err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(seed).Error
			if err != nil {
				errorMessage := err.Error()
				log.Printf("%s seeder failed with error: %s", name, errorMessage)
				continue
			}

			resetSequence(tx, seed)
		}

		log.Println("Seeding completed successfully.")
		return nil
	})
}

// resetSequence resets the auto-increment sequence in case of conflict or error
func resetSequence(tx *gorm.DB, seed interface{}) {
	if tx.Dialector.Name() == "postgres" {
		tableName := getTableName(seed)
		if tableName != "" {
			// Reset the sequence for PostgreSQL
			query := fmt.Sprintf(`
				SELECT setval(pg_get_serial_sequence('%s', 'id'), 
				COALESCE((SELECT MAX(id) FROM %s), 1))`, tableName, tableName)
			if err := tx.Exec(query).Error; err != nil {
				log.Printf("[WARNING] Failed to reset sequence for table %s: %s", tableName, err)
			}
		}
	}
}

// getTableName returns the table name from the seed's struct
func getTableName(seed interface{}) string {
	seedType := reflect.TypeOf(seed)
	if seedType.Kind() == reflect.Ptr {
		seedType = seedType.Elem()
	}
	if seedType.Kind() == reflect.Struct {
		return seedType.Name()
	}
	return ""
}

func dataSeeds() []interface{} {
	return []interface{}{
		models.UserSeed(),
	}
}
