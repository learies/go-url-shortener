package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/logger"
	"github.com/learies/go-url-shortener/internal/middlewares"
	"github.com/learies/go-url-shortener/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	// r.Use(middleware.Logger)
	r.Use(logger.WithLogging(log))
	r.Use(middlewares.GzipMiddleware)
	r.Mount("/", router.NewRouter(cfg))

	log.Info("Starting server...", slog.String("address", cfg.Address))
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Error("Server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
