package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupDB initializes and returns a gorm.DB connection
func SetupDB() (*gorm.DB, error) {
	//dsn := os.Getenv("DB_URL")
	dsn := "host=localhost user=postgres password=7968 dbname=goproj port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Adjust logging level here
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	return db, nil
}
