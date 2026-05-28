package storage

import (
	"context"
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

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
		seen_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS wishlist_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		album_name TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		UNIQUE(chat_id, album_name)
	);

	CREATE TABLE IF NOT EXISTS notification_settings (
		chat_id INTEGER PRIMARY KEY,
		mode TEXT NOT NULL
	);
`

func (s *Store) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, schema)
	return err
}
