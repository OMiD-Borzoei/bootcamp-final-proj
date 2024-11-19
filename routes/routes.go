package routes

import (
	"Project/config"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	db, _ := config.SetupDB()

	// Setup DL Routes
	SetupDLRoutes(router, db)

	// Setup SL Routes
	SetupSLRoutes(router, db)

	// Setup Voucher Routes
	SetupVoucherRoutes(router, db)
}
