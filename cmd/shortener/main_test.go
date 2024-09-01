package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		expectedCode   int
		expectedHeader string
		expectedBody   string
	}{
		{
			name:           "POST valid URL",
			method:         http.MethodPost,
			body:           "http://example.com",
			expectedCode:   http.StatusCreated,
			expectedHeader: "Content-Type",
			expectedBody:   "http://localhost:8080/",
		},
		{
			name:         "POST invalid URL",
			method:       http.MethodPost,
			body:         "invalid-url",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "GET existing short URL",
			method:       http.MethodGet,
			body:         "http://example.com",
			expectedCode: http.StatusTemporaryRedirect,
		},
		{
			name:         "GET non-existing short URL",
			method:       http.MethodGet,
			body:         "nonexisting",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Method Not Allowed",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/", mainHandler)

			// For POST method
			if tt.method == http.MethodPost {
				req, err := http.NewRequest(tt.method, "/", strings.NewReader(tt.body))
				assert.NoError(t, err)

				rec := httptest.NewRecorder()
				mux.ServeHTTP(rec, req)

				assert.Equal(t, tt.expectedCode, rec.Code)

				if tt.expectedHeader != "" {
					assert.Equal(t, tt.expectedBody[:len(tt.expectedBody)-1], rec.Body.String()[:len(tt.expectedBody)-1])
				}
			}

			// For GET method
			if tt.method == http.MethodGet {
				if tt.body != "nonexisting" {
					// First, create a shortened URL
					reqPost, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
					assert.NoError(t, err)

					recPost := httptest.NewRecorder()
					mux.ServeHTTP(recPost, reqPost)

					// Extract short URL from the POST response
					shortURL := strings.TrimPrefix(recPost.Body.String(), "http://localhost:8080/")
					req, err := http.NewRequest(tt.method, "/"+shortURL, nil)
					assert.NoError(t, err)

					rec := httptest.NewRecorder()
					mux.ServeHTTP(rec, req)

					assert.Equal(t, tt.expectedCode, rec.Code)
				} else {
					req, err := http.NewRequest(tt.method, "/"+tt.body, nil)
					assert.NoError(t, err)

					rec := httptest.NewRecorder()
					mux.ServeHTTP(rec, req)

					assert.Equal(t, tt.expectedCode, rec.Code)
				}
			}

			// For other methods (e.g., PUT)
			if tt.method == http.MethodPut {
				req, err := http.NewRequest(tt.method, "/", nil)
				assert.NoError(t, err)

				rec := httptest.NewRecorder()
				mux.ServeHTTP(rec, req)

				assert.Equal(t, tt.expectedCode, rec.Code)
			}
		})
	}
}
