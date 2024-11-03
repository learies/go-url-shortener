package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/learies/go-url-shortener/internal/config/logger"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

// Init initializes the database connection
func Init() {
	dbOnce.Do(func() {
		connStr := os.Getenv("DATABASE_URL")
		if connStr == "" {
			logger.Log.Error("DATABASE_URL is not set")
			return
		}

		var err error
		db, err = sql.Open("pgx", connStr)
		if err != nil {
			logger.Log.Error("Failed to open database connection", "error", err)
			return
		}

		if err := db.Ping(); err != nil {
			logger.Log.Error("Failed to ping database", "error", err)
			return
		}

		logger.Log.Info("Successfully connected to the database")
	})
}

// GetDB returns the database instance
func GetDB() (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return db, nil
}

// Close closes the database connection.
func Close() {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Log.Error("Failed to close database connection", "error", err)
		}
	}
}
