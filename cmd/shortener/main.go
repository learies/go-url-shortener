package main

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/learies/go-url-shortener/config"
)

var (
	urlMapping = make(map[string]string)
	mu         sync.Mutex
	cfg        *config.Config
)

func generateShortURL(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hash := hasher.Sum(nil)
	shortURL := base64.URLEncoding.EncodeToString(hash)[:8] // берем первые 8 байт для сокращения
	return shortURL
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read the request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		originalURL := string(body)
		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		mu.Lock()
		shortURL := generateShortURL(originalURL)
		urlMapping[shortURL] = originalURL
		mu.Unlock()

		shortenedURL := cfg.BaseURL + "/" + shortURL
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL))

	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/")
		mu.Lock()
		originalURL, exists := urlMapping[id]
		mu.Unlock()

		if !exists {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	cfg = config.ParseConfig()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", mainHandler)
	r.Get("/*", mainHandler)

	log.Printf("Starting server on %s...\n", cfg.Address)
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
