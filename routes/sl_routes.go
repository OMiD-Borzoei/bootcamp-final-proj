package routes

import (
	"Project/handlers"
	"Project/repositories"
	"Project/services"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupSLRoutes(router *mux.Router, db *gorm.DB) {
	slRepo := repositories.NewSLRepository(db)
	slService := services.NewSLService(slRepo)
	slHandler := handlers.NewSLHandler(slService)

	// Define SL API routes
	router.HandleFunc("/api/v1/sl", slHandler.GetAllSLs).Methods("GET")
	router.HandleFunc("/api/v1/sl", slHandler.CreateSL).Methods("POST")
	router.HandleFunc("/api/v1/sl/{id}", slHandler.GetSLByID).Methods("GET")
	router.HandleFunc("/api/v1/sl/{id}", slHandler.UpdateSL).Methods("PUT")
	router.HandleFunc("/api/v1/sl/{id}", slHandler.DeleteSL).Methods("DELETE")
}
