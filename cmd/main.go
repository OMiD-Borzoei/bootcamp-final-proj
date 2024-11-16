package main

import (
	"fmt"
	"log"

	"Project/config"
	"Project/models"
	"Project/repositories"

	"github.com/joho/godotenv"
	_ "gorm.io/driver/postgres"
)

func main() {
	godotenv.Load()

	// Set up database connection
	db, err := config.SetupDB()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	db.AutoMigrate(&models.DL{})

	dlRepo := repositories.NewDLRepository(db)

	// Example: Create a new DL record

	if err := dlRepo.Create("omid1", "reza1"); err != nil {
		fmt.Println("Error creating DL:", err)
	} else {
		fmt.Println("DL created successfully!")
	}

	dl, err := dlRepo.ReadByTitle("Test DL")
	if err != nil {
		fmt.Println("Error creating DL:", err)
	} else {
		fmt.Println(*dl)
	}

}
