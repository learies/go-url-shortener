package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/logger"
	"github.com/learies/go-url-shortener/internal/models"
	"github.com/learies/go-url-shortener/internal/shortener"
	"github.com/learies/go-url-shortener/internal/store"
)

func PostAPIHandler(store store.Store, cfg config.Config, urlShortener *shortener.URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Empty request body", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read the request body", http.StatusInternalServerError)
			return
		}

		var request models.Request
		if err = json.Unmarshal(body, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		originalURL := string(request.URL)
		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		shortURL := urlShortener.GenerateShortURL(originalURL)
		store.Set(shortURL, originalURL)
		shortenedURL := cfg.BaseURL + "/" + shortURL

		var response models.Response
		response.Result = shortenedURL

		result, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(result)
	}
}

func PostHandler(store store.Store, cfg config.Config, urlShortener *shortener.URLShortener) http.HandlerFunc {
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

func GetHandler(store store.Store) http.HandlerFunc {
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

// PingHandler проверяет доступность хранилища URL
func PingHandler(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := store.Ping(); err != nil {
			http.Error(w, "Store is not available", http.StatusInternalServerError)
			logger.Log.Error("Store ping failed", "error", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully connected to the store"))
	}
}
