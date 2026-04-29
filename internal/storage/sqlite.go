package storage

import (
	"context"
	"database/sql"
)

const schema = `
	CREATE TABLE IF NOT EXISTS subscribers (
		chat_id INTEGER PRIMARY KEY,
		created_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS albums (
		album_id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		artist TEXT NOT NULL,
		added_at DATETIME,
		seen_at DATETIME NOT NULL,
	);
`

func Init(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, schema)
	return err
}
