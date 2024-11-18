package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDB initializes and returns a gorm.DB connection
func SetupDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Silent), // Adjust logging level here
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	return db, nil
}
