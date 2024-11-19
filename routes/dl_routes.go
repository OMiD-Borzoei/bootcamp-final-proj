package routes

import (
	"Project/handlers"
	"Project/repositories"
	"Project/services"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupDLRoutes(router *mux.Router, db *gorm.DB) {
	repo := repositories.NewDLRepository(db)
	service := services.NewDLService(repo)
	handler := handlers.NewDLHandler(service)

	// Define DL API routes
	router.HandleFunc("/api/v1/dl", handler.GetAllDLs).Methods("GET")
	router.HandleFunc("/api/v1/dl/{id}", handler.GetDLByID).Methods("GET")
	router.HandleFunc("/api/v1/dl", handler.CreateDL).Methods("POST")
	router.HandleFunc("/api/v1/dl/{id}", handler.UpdateDL).Methods("PUT")
	router.HandleFunc("/api/v1/dl/{id}", handler.DeleteDL).Methods("DELETE")
}
