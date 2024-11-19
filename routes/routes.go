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

	// Voucher Handlers and Services
	vRepo := repositories.NewVoucherRepository(db)
	vService := services.NewVoucherService(vRepo)
	vHandler := handlers.NewVoucherHandler(vService)

	// Define Voucher API routes
	router.HandleFunc("/api/v1/v", vHandler.GetAllVouchers).Methods("GET")
	router.HandleFunc("/api/v1/v", vHandler.CreateVoucher).Methods("POST")
	router.HandleFunc("/api/v1/v/{id}", vHandler.GetVoucher).Methods("GET")
	router.HandleFunc("/api/v1/v/{id}", vHandler.UpdateVoucher).Methods("PUT")
	router.HandleFunc("/api/v1/v/{id}", vHandler.DeleteVoucher).Methods("DELETE")
}
