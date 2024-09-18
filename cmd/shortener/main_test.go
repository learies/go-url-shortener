package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/router"
	"github.com/learies/go-url-shortener/internal/shortener"
	"github.com/learies/go-url-shortener/internal/store"
)

func TestMainHandler(t *testing.T) {
	cfg := config.LoadConfig()
	cfg.BaseURL = "http://localhost:8080"

	store := store.NewURLStore()

	urlShortener := shortener.NewURLShortener()

	tests := []struct {
		name           string
		method         string
		body           string
		url            string
		expectedCode   int
		expectedHeader string
		expectedBody   string
	}{
		{
			name:           "POST valid URL",
			method:         http.MethodPost,
			body:           "http://example.com",
			url:            "/",
			expectedCode:   http.StatusCreated,
			expectedHeader: "Content-Type",
			expectedBody:   "http://localhost:8080/",
		},
		{
			name:         "POST invalid URL",
			method:       http.MethodPost,
			body:         "invalid-url",
			url:          "/",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "GET existing short URL",
			method:       http.MethodGet,
			body:         "http://example.com",
			url:          "/{id}",
			expectedCode: http.StatusTemporaryRedirect,
		},
		{
			name:         "GET non-existing short URL",
			method:       http.MethodGet,
			body:         "",
			url:          "/nonexisting",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Method Not Allowed",
			method:       http.MethodPut,
			body:         "",
			url:          "/",
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/", router.NewRouter(store, cfg, urlShortener))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Для метода POST
			if tt.method == http.MethodPost {
				req, err := http.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
				assert.NoError(t, err)

				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)

				assert.Equal(t, tt.expectedCode, rec.Code)

				if tt.expectedHeader != "" {
					assert.Equal(t, tt.expectedBody[:len(tt.expectedBody)-1], rec.Body.String()[:len(tt.expectedBody)-1])
				}
			}

			// Для метода GET
			if tt.method == http.MethodGet {
				if tt.body != "" {
					// Сначала создаем короткий URL
					reqPost, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
					assert.NoError(t, err)

					recPost := httptest.NewRecorder()
					r.ServeHTTP(recPost, reqPost)

					// Извлекаем короткий URL из ответа на POST
					shortURL := strings.TrimPrefix(recPost.Body.String(), cfg.BaseURL+"/")
					req, err := http.NewRequest(tt.method, "/"+shortURL, nil)
					assert.NoError(t, err)

					rec := httptest.NewRecorder()
					r.ServeHTTP(rec, req)

					assert.Equal(t, tt.expectedCode, rec.Code)
				} else {
					req, err := http.NewRequest(tt.method, tt.url, nil)
					assert.NoError(t, err)

					rec := httptest.NewRecorder()
					r.ServeHTTP(rec, req)

					assert.Equal(t, tt.expectedCode, rec.Code)
				}
			}

			// Для других методов (например, PUT)
			if tt.method != http.MethodPost && tt.method != http.MethodGet {
				req, err := http.NewRequest(tt.method, tt.url, nil)
				assert.NoError(t, err)

				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)

				assert.Equal(t, tt.expectedCode, rec.Code)
			}
		})
	}
}
