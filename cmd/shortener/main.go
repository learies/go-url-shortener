package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/", router.NewRouter(cfg))

	log.Printf("Starting server on %s...\n", cfg.Address)
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
