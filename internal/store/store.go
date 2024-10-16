package store

import (
	"database/sql"
	"encoding/json"
	"os"
	"sync"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/database"
	"github.com/learies/go-url-shortener/internal/logger"
)

// Store интерфейс для хранилища URL
type Store interface {
	Set(shortURL, originalURL string)
	Get(shortURL string) (string, bool)
	Ping() error
}

// URLStore хранение URL в файле
type URLStore struct {
	urlMapping map[string]string
	filePath   string
	mu         sync.Mutex
}

// URLMapping структура для JSON
type URLMapping struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Set сохраняет URL в память и файл
func (store *URLStore) Set(shortURL, originalURL string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.urlMapping[shortURL] = originalURL
	logger.Log.Info("Saving URLMapping", "shortURL", shortURL, "originalURL", originalURL)
	logger.Log.Info("Store", "filePath:", store.filePath)
	store.SaveToFile(store.filePath)
}

// Get получает URL из памяти или из файла
func (store *URLStore) Get(shortURL string) (string, bool) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.LoadFromFile(store.filePath)
	originalURL, exists := store.urlMapping[shortURL]
	return originalURL, exists
}

// SaveToFile сохраняет URL-маппинг в JSON файл
func (store *URLStore) SaveToFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for shortURL, originalURL := range store.urlMapping {
		logger.Log.Info("Encoding URLMapping", "shortURL", shortURL, "originalURL", originalURL)
		if err := encoder.Encode(URLMapping{ShortURL: shortURL, OriginalURL: originalURL}); err != nil {
			return err
		}
	}

	return nil
}

// LoadFromFile загружает URL-маппинг из JSON файла
func (store *URLStore) LoadFromFile(filePath string) error {
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
		store.urlMapping[urlMapping.ShortURL] = urlMapping.OriginalURL
	}

	return nil
}

// Ping проверяет доступность хранилища URL
func (store *URLStore) Ping() error {
	return nil // Файловое хранилище всегда доступно
}

// DBStore хранение URL в базе данных
type DBStore struct {
	db *sql.DB
}

// Set сохраняет URL в базу данных
func (ds *DBStore) Set(shortURL, originalURL string) {
	_, err := ds.db.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (short_url) DO UPDATE SET original_url = $2", shortURL, originalURL)
	if err != nil {
		logger.Log.Error("Failed to set URL mapping in database", "error", err)
	}
}

// Get получает URL из базы данных
func (ds *DBStore) Get(shortURL string) (string, bool) {
	var originalURL string
	err := ds.db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
		logger.Log.Error("Failed to get URL mapping from database", "error", err)
		return "", false
	}
	return originalURL, true
}

// Ping проверяет доступность базы данных
func (ds *DBStore) Ping() error {
	return ds.db.Ping()
}

// NewStore создаёт новое хранилище URL
func NewStore(cfg config.Config) (Store, error) {
	if cfg.DatabaseDSN != "" {
		// Подключение к базе данных
		db, err := database.Connect(cfg.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		store := &DBStore{db: db}
		return store, nil
	}

	// Используем файловое хранилище
	store := &URLStore{
		urlMapping: make(map[string]string),
		filePath:   cfg.FileStoragePath,
	}

	if store.filePath == "" {
		// Создание временного файла
		tmpFile, err := os.CreateTemp("/tmp", "urlstore-*.json")
		if err != nil {
			return nil, err
		}
		defer tmpFile.Close()
		store.filePath = tmpFile.Name()
	}
	store.LoadFromFile(store.filePath)
	return store, nil
}
