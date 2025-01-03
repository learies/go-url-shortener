package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/contextutils"
	"github.com/learies/go-url-shortener/internal/logger"
	"github.com/learies/go-url-shortener/internal/models"
	"github.com/learies/go-url-shortener/internal/shortener"
	"github.com/learies/go-url-shortener/internal/store"
	"github.com/learies/go-url-shortener/internal/worker"
)

func PostAPIHandler(store store.Store, cfg config.Config, urlShortener *shortener.URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

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

		var response models.Response
		response.Result = cfg.BaseURL + "/" + shortURL

		result, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		// Получим userID из контекста
		userID, ok := contextutils.GetUserID(ctx)
		if !ok {
			http.Error(w, "UserID not found in context", http.StatusUnauthorized)
			return
		}

		err = store.Set(ctx, shortURL, originalURL, userID)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to store URL: %v", err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write(result)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(result)
	}
}

func PostAPIBatchHandler(store store.Store, cfg config.Config, urlShortener *shortener.URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		// Получим userID из контекста
		userID, ok := contextutils.GetUserID(ctx)
		if !ok {
			http.Error(w, "UserID not found in context", http.StatusUnauthorized)
			return
		}

		if r.Body == nil {
			http.Error(w, "Empty request body", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read the request body", http.StatusInternalServerError)
			return
		}

		var requests []models.BatchURLRequest
		if err = json.Unmarshal(body, &requests); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var responses []models.BatchURLResponse
		var batchWrites []models.BatchURLWrite
		for _, request := range requests {
			shortURL := urlShortener.GenerateShortURL(request.OriginalURL)
			responses = append(responses, models.BatchURLResponse{
				CorrelationID: request.CorrelationID,
				ShortURL:      cfg.BaseURL + "/" + shortURL,
			})
			batchWrites = append(batchWrites, models.BatchURLWrite{
				CorrelationID: request.CorrelationID,
				ShortURL:      shortURL,
				OriginalURL:   request.OriginalURL,
				UserID:        userID,
			})
		}

		store.SetBatch(ctx, batchWrites)

		result, err := json.Marshal(responses)
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
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

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

		// Получим userID из контекста
		userID, ok := contextutils.GetUserID(ctx)
		if !ok {
			http.Error(w, "UserID not found in context", http.StatusUnauthorized)
			return
		}

		shortURL := urlShortener.GenerateShortURL(originalURL)
		shortenedURL := cfg.BaseURL + "/" + shortURL

		err = store.Set(ctx, shortURL, originalURL, userID)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to store URL: %v", err))
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(shortenedURL))
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL))
	}
}

func GetHandler(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		shortURL := strings.TrimPrefix(r.URL.Path, "/")

		s, exists := store.Get(ctx, shortURL)
		if !exists {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		if s.DeletedFlag {
			http.Error(w, "URL is deleted", http.StatusGone)
			return
		}

		w.Header().Set("Location", s.OriginalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func GetAPIUserURLsHandler(store store.Store, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		// Получим userID из контекста
		userID, ok := contextutils.GetUserID(ctx)
		if !ok {
			http.Error(w, "UserID not found in context", http.StatusUnauthorized)
			return
		}

		urls, ok := store.GetUserUrls(ctx, userID)
		if !ok {
			http.Error(w, "URLs not found", http.StatusNotFound)
			return
		}

		if len(urls) == 0 {
			w.WriteHeader(http.StatusUnauthorized) // Переделать на 204
			return
		}

		// Модифицируем URL для возврата
		modifiedUrls := make([]models.URL, len(urls))

		// Для каждого URL добавим BaseURL
		for i, url := range urls {
			modifiedUrls[i] = models.URL{
				ShortURL:    cfg.BaseURL + "/" + url.ShortURL,
				OriginalURL: url.OriginalURL,
			}
		}

		result, err := json.Marshal(modifiedUrls)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func DeleteUserUrlsHandler(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		userID, ok := contextutils.GetUserID(ctx)
		if !ok {
			http.Error(w, "UserID not found in context", http.StatusUnauthorized)
			return
		}

		var shortURLs models.ShortURLs
		if err := json.NewDecoder(r.Body).Decode(&shortURLs.ShortURLs); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		for _, shortURL := range shortURLs.ShortURLs {
			store.DeleteUserUrls(ctx, worker.GenerateShortURL(models.UserURL{
				UserID:   userID,
				ShortURL: shortURL,
			}))
		}

		w.WriteHeader(http.StatusAccepted)
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
