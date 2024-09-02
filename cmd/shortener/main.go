package main

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	urlMapping = make(map[string]string)
	mu         sync.Mutex
)

func generateShortURL(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hash := hasher.Sum(nil)
	shortURL := base64.URLEncoding.EncodeToString(hash)[:8] // берем первое 8 байт для сокращения
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

		shortenedURL := "http://localhost:8080/" + shortURL
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)

	log.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
