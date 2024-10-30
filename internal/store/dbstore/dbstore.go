package dbstore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/learies/go-url-shortener/internal/logger"
	"github.com/learies/go-url-shortener/internal/models"
)

// DBStore хранение URL в базе данных
type DBStore struct {
	DB *sql.DB
}

// Set сохраняет URL в базу данных
func (ds *DBStore) Set(ctx context.Context, shortURL, originalURL, userID string) error {
	id := uuid.New()

	query := `
	INSERT INTO urls (id, short_url, original_url, user_id)
	VALUES ($1, $2, $3, $4)`
	// ON CONFLICT (short_url) DO UPDATE SET original_url = EXCLUDED.original_url;`

	_, err := ds.DB.ExecContext(ctx, query, id, shortURL, originalURL, userID)
	if err != nil {
		return err
	}
	return nil
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

// SetBatch сохраняет пакет URL в базе данных
func (ds *DBStore) SetBatch(ctx context.Context, urls []models.BatchURLWrite) {
	if len(urls) == 0 {
		return
	}

	var (
		queryValues string
		args        []interface{}
	)

	query := `
    INSERT INTO urls (id, short_url, original_url, user_id)
    VALUES %s
    ON CONFLICT (short_url) DO UPDATE SET original_url = EXCLUDED.original_url;`

	for i, response := range urls {
		// Создание плейсхолдера для каждой строки
		// ($1, $2, $3), ($4, $5, $6), ...
		placeholder := fmt.Sprintf("($%d, $%d, $%d, $%d)", (i*4)+1, (i*4)+2, (i*4)+3, (i*4)+4)
		queryValues += placeholder
		if i < len(urls)-1 {
			queryValues += ", "
		}

		args = append(args, response.CorrelationID, response.ShortURL, response.OriginalURL, response.UserID)
	}

	fullQuery := fmt.Sprintf(query, queryValues)
	_, err := ds.DB.ExecContext(ctx, fullQuery, args...)
	if err != nil {
		logger.Log.Error("Failed to set URL mapping in database", "error", err)
	}
}

// GetBatch получает пакет URL из базы данных
func (ds *DBStore) GetUserUrls(ctx context.Context, userID string) ([]models.Url, bool) {
	var urls []models.Url

	rows, err := ds.DB.QueryContext(ctx, "SELECT short_url, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		logger.Log.Error("Failed to get user URLs from database", "error", err)
		return nil, false
	}
	defer rows.Close()

	for rows.Next() {
		var url models.Url
		err := rows.Scan(&url.ShortURL, &url.OriginalURL)
		if err != nil {
			logger.Log.Error("Failed to scan user URLs from database", "error", err)
			return nil, false
		}
		urls = append(urls, url)
	}

	return urls, true
}

// Ping проверяет доступность базы данных
func (ds *DBStore) Ping() error {
	return ds.DB.Ping()
}
