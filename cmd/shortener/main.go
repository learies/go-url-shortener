package main

import (
	"log"
	"net/http"
	"os"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/logger"
	"github.com/learies/go-url-shortener/internal/router"
)

func main() {
	err := logger.Initialize("info")
	if err != nil {
		log.Println("Error initializing logger", "err", err)
		return
	}

	cfg := config.LoadConfig()

	r := router.NewRouter(cfg)

	logger.Log.Info("Starting server", "address", cfg.Address)
	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		logger.Log.Error("Error starting server", "err", err)
		os.Exit(1)
	}
}
