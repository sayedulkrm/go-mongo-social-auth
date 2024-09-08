package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/sayedulkrm/go-mongo-social-auth/config"
	"github.com/sayedulkrm/go-mongo-social-auth/lib"
	"github.com/sayedulkrm/go-mongo-social-auth/middlewares"
	"github.com/sayedulkrm/go-mongo-social-auth/routes"
	"github.com/sirupsen/logrus"
)

func main() {

	// Configure Logger
	lib.ConfigureLogger()

	// routes calls
	root := routes.SetupRoutes()

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		// Debug:            true,

	}

	cors := cors.New(corsOptions)

	handler := cors.Handler(root)
	errorHandler := middlewares.ErrorMiddleware(handler)

	startServer(errorHandler)

}

func startServer(handler http.Handler) {
	// Start Server
	err := godotenv.Load(".env")

	if err != nil {
		logrus.Fatalf("Error loading .env file")
	}

	// Db Connect
	config.DBInstance()

	port := os.Getenv("PORT")

	if port == "" {
		port = ":8000"
	}

	logrus.Warn("Server running on", port)

	if err := http.ListenAndServe(port, handler); err != nil {
		logrus.Fatal("Failed to start server", err)
	}

}
