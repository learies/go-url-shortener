package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/handlers"
	"github.com/learies/go-url-shortener/internal/shortener"
	"github.com/learies/go-url-shortener/internal/store"
)

func NewRouter(cfg config.Config) http.Handler {
	store := store.NewURLStore()
	urlShortener := shortener.NewURLShortener()

	r := chi.NewRouter()
	r.Post("/", handlers.PostHandler(store, cfg, urlShortener))
	r.Post("/api/shorten", handlers.PostAPIHandler(store, cfg, urlShortener))
	r.Get("/*", handlers.GetHandler(store))

	r.MethodNotAllowed(methodNotAllowedHandler)
	return r
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
