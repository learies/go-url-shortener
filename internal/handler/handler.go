package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/learies/go-url-shortener/internal/config/logger"
	"github.com/learies/go-url-shortener/internal/db"
)

func HelloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}
}

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := db.GetDB()
		if err != nil {
			http.Error(w, "database not initialized", http.StatusInternalServerError)
			logger.Log.Error("Database not initialized", "error", err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				http.Error(w, "Database ping timed out", http.StatusRequestTimeout)
				logger.Log.Error("Database ping timed out", "error", err)
				return
			}
			http.Error(w, "Failed to ping database", http.StatusInternalServerError)
			logger.Log.Error("Failed to ping database", "error", err)
			return
		}

		w.Write([]byte("Successfully connected to the database"))
	}
}
