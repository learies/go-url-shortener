package main

import (
	"log"
	"net/http"
	"os"

	"github.com/learies/go-url-shortener/internal/config/logger"
	"github.com/learies/go-url-shortener/internal/db"
	"github.com/learies/go-url-shortener/internal/router"
)

func main() {
	// Initialize the logger
	err := logger.InitLogger("info")
	if err != nil {
		log.Fatalf("could not initialize logger: %v", err)
	}

	// Initialize the database
	db.Init()
	defer db.Close()

	r := router.InitRouter()

	// Start the server
	logger.Log.Info("Starting server", "port", 8080)
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Log.Error("could not start server", "error", err)
		os.Exit(1)
	}
}
