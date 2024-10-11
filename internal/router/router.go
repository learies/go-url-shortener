package router

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/handlers"
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
		log.Fatalf("Error opening database connection: %v", err)
		os.Exit(1)
	}

	r := chi.NewRouter()
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
