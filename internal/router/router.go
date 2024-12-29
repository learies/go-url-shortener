package router

import (
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

func NewRouter(cfg config.Config) http.Handler {
	store, err := store.NewStore(cfg)
	if err != nil {
		logger.Log.Error("Error creating store", "err", err)
		os.Exit(1)
	}

	urlShortener := shortener.NewURLShortener()

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(internalMiddleware.WithLogging)
	r.Use(internalMiddleware.GzipMiddleware)
	r.Use(internalMiddleware.JWTMiddleware)

	r.Post("/", handlers.PostHandler(store, cfg, urlShortener))
	r.Post("/api/shorten", handlers.PostAPIHandler(store, cfg, urlShortener))
	r.Post("/api/shorten/batch", handlers.PostAPIBatchHandler(store, cfg, urlShortener))
	r.Get("/api/user/urls", handlers.GetAPIUserURLsHandler(store, cfg))
	r.Delete("/api/user/urls", handlers.DeleteUserUrlsHandler(store))
	r.Get("/*", handlers.GetHandler(store))
	r.Get("/ping", handlers.PingHandler(store))

	r.MethodNotAllowed(methodNotAllowedHandler)

	return r
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
