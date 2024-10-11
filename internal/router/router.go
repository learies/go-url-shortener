package router

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/handlers"
	"github.com/learies/go-url-shortener/internal/logger"
	"github.com/learies/go-url-shortener/internal/middlewares"
	"github.com/learies/go-url-shortener/internal/shortener"
	"github.com/learies/go-url-shortener/internal/store"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func connectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewRouter(cfg config.Config) http.Handler {
	store := store.NewURLStore(cfg.FileStoragePath)
	urlShortener := shortener.NewURLShortener()

	db, err := connectToDB(cfg.DatabaseDSN)
	if err != nil {
		slog.Error("Error opening database connection", "err", err)
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	// r.Use(middleware.Logger)
	r.Use(logger.WithLogging(slog.New(slog.NewJSONHandler(os.Stdout, nil))))
	r.Use(middlewares.GzipMiddleware)

	r.Post("/", handlers.PostHandler(store, cfg, urlShortener))
	r.Post("/api/shorten", handlers.PostAPIHandler(store, cfg, urlShortener))
	r.Get("/*", handlers.GetHandler(store))
	r.Get("/ping", handlers.PingHandler(db))

	r.MethodNotAllowed(methodNotAllowedHandler)

	return r
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
