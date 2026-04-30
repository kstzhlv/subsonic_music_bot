package storage

import (
	"context"
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{}
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
		seen_at DATETIME NOT NULL,
	);
`

func (s *Store) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, schema)
	return err
}

func (s *Store) AddSubscriber(ctx context.Context, chatID int64) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT OR IGNORE INTO subscribers (chat_id, created_at) VALUES (?, ?)`,
		chatID,
		time.Now().UTC(),
	)

	return err
}

func (s *Store) RemoveSubscriber(ctx context.Context, chatID int64) error {
	_, err := s.db.ExecContext(
		ctx,
		`DELETE FROM subscribers WHERE chat_id = ?`,
		chatId,
	)

	return err
}
