package filestore

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/learies/go-url-shortener/internal/logger"
	"github.com/learies/go-url-shortener/internal/models"
)

// URLStore хранение URL в файле
type FileStore struct {
	URLMapping map[string]string
	FilePath   string
	mu         sync.Mutex
}

// URLMapping структура для JSON
type URLMapping struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Set сохраняет URL в память и файл
func (store *FileStore) Set(ctx context.Context, shortURL, originalURL string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.URLMapping[shortURL] = originalURL
	logger.Log.Info("Saving URLMapping", "shortURL", shortURL, "originalURL", originalURL)
	logger.Log.Info("Store", "filePath:", store.FilePath)
	store.SaveToFile(store.FilePath)
	return nil
}

// Get получает URL из памяти или из файла
func (store *FileStore) Get(ctx context.Context, shortURL string) (string, bool) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.LoadFromFile(store.FilePath)
	originalURL, exists := store.URLMapping[shortURL]
	return originalURL, exists
}

// SetBatch сохраняет URL в память и файл
func (store *FileStore) SetBatch(ctx context.Context, shortURLS []models.BatchURLWrite) {
	store.mu.Lock()
	defer store.mu.Unlock()
	for _, urlMapping := range shortURLS {
		store.URLMapping[urlMapping.ShortURL] = urlMapping.OriginalURL
		logger.Log.Info("Saving URLMapping", "shortURL", urlMapping.ShortURL, "originalURL", urlMapping.OriginalURL)
	}
	store.SaveToFile(store.FilePath)
}

// SaveToFile сохраняет URL-маппинг в JSON файл
func (store *FileStore) SaveToFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for shortURL, originalURL := range store.URLMapping {
		logger.Log.Info("Encoding URLMapping", "shortURL", shortURL, "originalURL", originalURL)
		if err := encoder.Encode(URLMapping{ShortURL: shortURL, OriginalURL: originalURL}); err != nil {
			return err
		}
	}

	return nil
}

// LoadFromFile загружает URL-маппинг из JSON файла
func (store *FileStore) LoadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
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
		store.URLMapping[urlMapping.ShortURL] = urlMapping.OriginalURL
	}

	return nil
}

// Ping проверяет доступность хранилища URL
func (store *FileStore) Ping() error {
	err := errors.New("unable to access the store")
	return err
}
