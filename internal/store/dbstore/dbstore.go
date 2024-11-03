package dbstore

import (
	"context"
	"database/sql"

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
	tx, err := ds.DB.Begin()
	if err != nil {
		logger.Log.Error("Failed to start transaction", "error", err)
		return
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO urls (id, short_url, original_url, user_id) VALUES ($1, $2, $3, $4)")
	if err != nil {
		logger.Log.Error("Failed to prepare statement", "error", err)
		tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, url := range urls {
		_, err := stmt.ExecContext(ctx, url.CorrelationID, url.ShortURL, url.OriginalURL, url.UserID)
		if err != nil {
			logger.Log.Error("Failed to insert URL", "error", err)
			tx.Rollback()
			logger.Log.Info("Transaction rolled back")
		}
	}
	err = tx.Commit()
	if err != nil {
		logger.Log.Error("Failed to commit transaction", "error", err)
	}
}

// GetBatch получает пакет URL из базы данных
func (ds *DBStore) GetUserUrls(ctx context.Context, userID string) ([]models.URL, bool) {
	var urls []models.URL

	rows, err := ds.DB.QueryContext(ctx, "SELECT short_url, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		logger.Log.Error("Failed to get user URLs from database", "error", err)
		return nil, false
	}
	defer rows.Close()

	for rows.Next() {
		var url models.URL
		err := rows.Scan(&url.ShortURL, &url.OriginalURL)
		if err != nil {
			logger.Log.Error("Failed to scan user URLs from database", "error", err)
			return nil, false
		}
		urls = append(urls, url)
	}

	// Check for any error after closing the loop
	if err = rows.Err(); err != nil {
		logger.Log.Error("Failed during rows iteration", "error", err)
		return nil, false
	}

	return urls, true
}

// DeleteUserUrls устанавливает флаг is_deleted в true для URL, принадлежащих пользователю.
func (ds *DBStore) DeleteUserUrls(ctx context.Context, userID string, shortURLs []string) error {
	query := "UPDATE urls SET is_deleted = TRUE WHERE user_id = $1 AND short_url = ANY($2)"
	_, err := ds.DB.ExecContext(ctx, query, userID, shortURLs)
	if err != nil {
		return err
	}
	return nil
}

// Ping проверяет доступность базы данных
func (ds *DBStore) Ping() error {
	return ds.DB.Ping()
}
