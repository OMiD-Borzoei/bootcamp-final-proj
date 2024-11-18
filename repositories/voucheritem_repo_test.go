package repositories

import (
	"Project/config"
	"Project/models"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestAllVImethods(t *testing.T) {
	// Set up database connection
	godotenv.Load()
	db, err := config.SetupDB()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	repo := NewDLRepository(db)

	repo.db.AutoMigrate(&models.DL{})

}
