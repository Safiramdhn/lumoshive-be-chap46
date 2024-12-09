package database

import (
	"fmt"
	"golang-chap46/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(config config.Configuration) (*gorm.DB, error) {
	// Format connection string
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s TimeZone=%s",
		config.DBConfig.DBUsername, config.DBConfig.DBPassword, config.DBConfig.DBName, config.DBConfig.DBHost, config.DBConfig.DBTimeZone)

	// Setup logger for GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	// Open a connection to the PostgreSQL databas
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	// Convert to *sql.DB for setting connection options
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool options
	sqlDB.SetConnMaxIdleTime(time.Duration(config.DBConfig.DBMaxIdleTime) * time.Minute)
	sqlDB.SetConnMaxLifetime(time.Duration(config.DBConfig.DBMaxLifeTime) * time.Hour)
	sqlDB.SetMaxIdleConns(config.DBConfig.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.DBConfig.DBMaxOpenConns)

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("ERROR: Database connection failed: %v\n", err)
		return nil, fmt.Errorf("ERROR: unable to connect to the database: %v", err)
	}
	log.Println("INFO: Database connected successfully!")

	// Migration tabel form struct
	log.Println("Starting migration...")
	err = Migrate(db)
	if err != nil {
		log.Fatalf("ERROR: unable to migrate database: %v", err)
	}
	log.Println("Migration completed successfully.")

	// running seeder
	log.Println("Starting seeding...")
	err = SeedAll(db)
	if err != nil {
		log.Fatalf("ERROR: failed seedingAll, message: %s", err.Error())
	}

	return db, nil
}
