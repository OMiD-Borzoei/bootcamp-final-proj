package routes

import (
	"Project/handlers"
	"Project/repositories"
	"Project/services"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupVoucherRoutes(router *mux.Router, db *gorm.DB) {
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
