package router

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/handlers"
	"github.com/learies/go-url-shortener/internal/logger"
	internalMiddleware "github.com/learies/go-url-shortener/internal/middleware"
	"github.com/learies/go-url-shortener/internal/shortener"
	"github.com/learies/go-url-shortener/internal/store"
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
		logger.Log.Error("Error opening database connection", "err", err)
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	// r.Use(middleware.Logger)
	r.Use(internalMiddleware.WithLogging)
	r.Use(internalMiddleware.GzipMiddleware)

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
