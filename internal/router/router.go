package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/handlers"
	"github.com/learies/go-url-shortener/internal/store"
)

func NewRouter(store *store.URLStore, cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Post("/", handlers.PostHandler(store, cfg))
	r.Get("/*", handlers.GetHandler(store))

	r.MethodNotAllowed(methodNotAllowedHandler)
	return r
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
