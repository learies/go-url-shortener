package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	r := router.NewRouter(cfg)

	slog.Info("Starting server...", slog.String("address", cfg.Address))
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		slog.Error("Server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
