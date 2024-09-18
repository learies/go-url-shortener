package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/shortener"
	"github.com/learies/go-url-shortener/internal/store"
)

func PostHandler(store *store.URLStore, cfg config.Config, urlShortener *shortener.URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		shortURL := urlShortener.GenerateShortURL(originalURL)
		store.Set(shortURL, originalURL)

		shortenedURL := cfg.BaseURL + "/" + shortURL
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL))
	}
}

func GetHandler(store *store.URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/")

		originalURL, exists := store.Get(id)
		if !exists {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
