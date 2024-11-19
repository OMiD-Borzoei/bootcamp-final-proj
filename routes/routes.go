package routes

import (
	"Project/config"
	"Project/handlers"
	"Project/repositories"
	"Project/services"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	db, _ := config.SetupDB()

	// DL Handlers and Services
	repo := repositories.NewDLRepository(db)
	service := services.NewDLService(repo)
	handler := handlers.NewDLHandler(service)

	// Define DL API routes
	router.HandleFunc("/api/v1/dl", handler.GetAllDLs).Methods("GET")
	router.HandleFunc("/api/v1/dl/{id}", handler.GetDLByID).Methods("GET")
	router.HandleFunc("/api/v1/dl", handler.CreateDL).Methods("POST")
	router.HandleFunc("/api/v1/dl/{id}", handler.UpdateDL).Methods("PUT")
	router.HandleFunc("/api/v1/dl/{id}", handler.DeleteDL).Methods("DELETE")

	// SL Handlers and Services
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
