package store

import "sync"

type URLStore struct {
	urlMapping map[string]string
	mu         sync.Mutex
}

func NewURLStore() *URLStore {
	return &URLStore{
		urlMapping: make(map[string]string),
	}
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
