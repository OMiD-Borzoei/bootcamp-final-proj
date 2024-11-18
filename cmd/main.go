package main

import (
	"log"

	"Project/config"

	"github.com/joho/godotenv"
	_ "gorm.io/driver/postgres"
)

func main() {
	godotenv.Load()

	// Set up database connection
	_, err := config.SetupDB()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
}
