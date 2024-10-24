package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/logger"
	"github.com/learies/go-url-shortener/internal/models"
	"github.com/learies/go-url-shortener/internal/router"
)

func TestMainHandler(t *testing.T) {

	cfg := config.LoadConfig()

	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Println("Error initializing logger", "err", err)
		return
	}

	cfg.BaseURL = "http://localhost:8080"

	r := router.NewRouter(cfg)

	t.Run("POST valid URL", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "text/plain")
		assert.Contains(t, rec.Body.String(), "http://localhost:8080/")
	})

	t.Run("POST invalid URL", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("invalid-url"))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GET existing short URL", func(t *testing.T) {
		// Сначала создаем короткий URL
		reqPost, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
		assert.NoError(t, err)

		recPost := httptest.NewRecorder()
		r.ServeHTTP(recPost, reqPost)

		// Извлекаем короткий URL из ответа на POST
		shortURL := strings.TrimPrefix(recPost.Body.String(), cfg.BaseURL+"/")

		req, err := http.NewRequest(http.MethodGet, "/"+shortURL, nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	})

	t.Run("GET non-existing short URL", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/nonexisting", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, "/", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("POST /api/shorten valid URL", func(t *testing.T) {
		requestBody, _ := json.Marshal(models.Request{URL: "http://example.com"})
		req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(requestBody))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	})

	t.Run("POST /api/shorten invalid URL", func(t *testing.T) {
		requestBody, _ := json.Marshal(models.Request{URL: "invalid-url"})
		req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(requestBody))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("POST /api/shorten invalid JSON", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader("{invalid-json}"))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("POST /api/shorten empty body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/shorten", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
