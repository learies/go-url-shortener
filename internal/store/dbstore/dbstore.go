package dbstore

import (
	"context"
	"database/sql"

	"github.com/learies/go-url-shortener/internal/logger"
)

// DBStore хранение URL в базе данных
type DBStore struct {
	DB *sql.DB
}

// Set сохраняет URL в базу данных
func (ds *DBStore) Set(ctx context.Context, shortURL, originalURL string) {
	_, err := ds.DB.ExecContext(ctx, "INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (short_url) DO UPDATE SET original_url = $2", shortURL, originalURL)
	if err != nil {
		logger.Log.Error("Failed to set URL mapping in database", "error", err)
	}
}

// Get получает URL из базы данных
func (ds *DBStore) Get(ctx context.Context, shortURL string) (string, bool) {
	var originalURL string
	err := ds.DB.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
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
	return ds.DB.Ping()
}
