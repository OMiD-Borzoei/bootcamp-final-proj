package main

import (
	"Project/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize the router
	router := mux.NewRouter()

	// Setup the routes
	routes.SetupRoutes(router)

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins or specify a list of allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Use the CORS handler
	handler := corsHandler.Handler(router)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", handler))
}
