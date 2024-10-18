package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := initialize(db); err != nil {
		return nil, err
	}

	return db, nil
}

func initialize(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id UUID PRIMARY KEY,
	    short_url VARCHAR(8) NOT NULL UNIQUE,
	    original_url TEXT NOT NULL UNIQUE,
	    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
