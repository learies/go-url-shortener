package store

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/learies/go-url-shortener/internal/logger"
)

type URLStore struct {
	urlMapping map[string]string
	filePath   string
	mu         sync.Mutex
}

type URLMapping struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewURLStore(filePath string) *URLStore {
	store := &URLStore{
		urlMapping: make(map[string]string),
		filePath:   filePath,
	}
	if store.filePath == "" {
		tmpFile, err := os.CreateTemp("/tmp", "urlstore-*.json")
		if err != nil {
			panic(err)
		}
		defer tmpFile.Close()
		store.filePath = tmpFile.Name()
	}
	store.LoadFromFile(filePath)
	return store
}

func (store *URLStore) Set(shortURL, originalURL string) {
	store.urlMapping[shortURL] = originalURL
	logger.Log.Info("Saving URLMapping", "shortURL", shortURL, "originalURL", originalURL)
	logger.Log.Info("Store", "filePath:", store.filePath)
	store.SaveToFile(store.filePath)
}

func (store *URLStore) Get(shortURL string) (string, bool) {
	store.LoadFromFile(store.filePath)
	originalURL, exists := store.urlMapping[shortURL]
	return originalURL, exists
}

func (store *URLStore) SaveToFile(filePath string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	// Открываем файл для записи в конце файла или создаем новый
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Кодируем URL-маппинг в JSON и записываем в файл
	encoder := json.NewEncoder(file)
	for shortURL, originalURL := range store.urlMapping {
		logger.Log.Info("Encoding URLMapping", "shortURL", shortURL, "originalURL", originalURL)
		if err := encoder.Encode(URLMapping{ShortURL: shortURL, OriginalURL: originalURL}); err != nil {
			return err
		}
	}

	return nil
}

func (store *URLStore) LoadFromFile(filePath string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	// Открываем файл для чтения
	// Если файл не существует, возвращаем nil
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	// Декодируем JSON из файла
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
