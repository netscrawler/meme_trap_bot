package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func NewStorage(ctx context.Context, dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err := db.ExecContext(ctx, imagesTable); err != nil {
		return nil, fmt.Errorf("failed to create images table: %w", err)
	}

	return db, nil
}
