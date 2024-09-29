package store

import (
	"encoding/json"
	"os"
	"sync"
)

type URLStore struct {
	urlMapping map[string]string
	mu         sync.Mutex
}

type URLMapping struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewURLStore(filePath string) *URLStore {
	store := &URLStore{
		urlMapping: make(map[string]string),
	}
	store.LoadFromFile(filePath)
	return store
}

func (store *URLStore) Set(shortURL, originalURL string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.urlMapping[shortURL] = originalURL
}

func (store *URLStore) Get(shortURL string) (string, bool) {
	store.mu.Lock()
	defer store.mu.Unlock()
	originalURL, exists := store.urlMapping[shortURL]
	return originalURL, exists
}

func (store *URLStore) SaveToFile(filePath string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for shortURL, originalURL := range store.urlMapping {
		if err := encoder.Encode(URLMapping{ShortURL: shortURL, OriginalURL: originalURL}); err != nil {
			return err
		}
	}

	return nil
}

func (store *URLStore) LoadFromFile(filePath string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File does not exist, nothing to load
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var urlMapping URLMapping
		if err := decoder.Decode(&urlMapping); err != nil {
			break
		}
		store.urlMapping[urlMapping.ShortURL] = urlMapping.OriginalURL
	}

	return nil
}
