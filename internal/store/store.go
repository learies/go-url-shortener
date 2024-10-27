package store

import (
	"context"
	"os"

	"github.com/learies/go-url-shortener/config"
	"github.com/learies/go-url-shortener/internal/database"
	"github.com/learies/go-url-shortener/internal/models"
	"github.com/learies/go-url-shortener/internal/store/dbstore"
	"github.com/learies/go-url-shortener/internal/store/filestore"
)

// Store интерфейс для хранилища URL
type Store interface {
	Set(ctx context.Context, shortURL, originalURL, userID string) error
	Get(ctx context.Context, shortURL string) (string, bool)
	SetBatch(ctx context.Context, shortURLS []models.BatchURLWrite)
	Ping() error
}

// NewStore создаёт новое хранилище URL
func NewStore(cfg config.Config) (Store, error) {
	if cfg.DatabaseDSN != "" {
		// Подключение к базе данных
		db, err := database.Connect(cfg.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		store := &dbstore.DBStore{DB: db}
		return store, nil
	}

	// Используем файловое хранилище
	store := &filestore.FileStore{
		URLMapping: make(map[string]string),
		FilePath:   cfg.FileStoragePath,
	}

	if store.FilePath == "" {
		// Создание временного файла
		tmpFile, err := os.CreateTemp("/tmp", "urlstore-*.json")
		if err != nil {
			return nil, err
		}
		defer tmpFile.Close()
		store.FilePath = tmpFile.Name()
	}
	store.LoadFromFile(store.FilePath)
	return store, nil
}
