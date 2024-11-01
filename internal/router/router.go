package router

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/learies/go-url-shortener/internal/handler"
	internalMiddleware "github.com/learies/go-url-shortener/internal/middleware"
)

func InitRouter() http.Handler {
	r := chi.NewRouter()

	// Use chi's built-in middleware for logging and recovery
	r.Use(middleware.Recoverer)

	// Use our custom logger middleware
	r.Use(internalMiddleware.MiddlewareLogger)

	// Define routes
	r.Get("/", handler.HelloHandler())

	return r
}
