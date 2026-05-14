package storage

import (
	"context"
	"database/sql"
	"time"

	"example.com/subsonic_bot/internal/domain"
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
		mode TEXT NOT NULL,
	)
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
		chatID,
	)

	return err
}

func (s *Store) ListSubscribers(ctx context.Context) ([]int64, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT chat_id FROM subscribers`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribers []int64
	for rows.Next() {
		var chatID int64
		if err := rows.Scan(&chatID); err != nil {
			return nil, err
		}
		subscribers = append(subscribers, chatID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscribers, nil
}

func (s *Store) AlbumIsSeen(
	ctx context.Context,
	albumID string,
) (bool, error) {
	var exists int

	err := s.db.QueryRowContext(
		ctx,
		`SELECT 1 FROM albums WHERE album_id = ? LIMIT 1`,
		albumID,
	).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Store) SaveSeenAlbum(
	ctx context.Context,
	album domain.Album,
) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT OR IGNORE INTO albums
		(album_id, title, artist, added_at, seen_at)
		VALUES (?, ?, ?, ?, ?)`,
		album.ID,
		album.Title,
		album.Artist,
		nullableTime(album.AddedAt),
		time.Now().UTC(),
	)

	return err
}

func (s *Store) CountSeenAlbums(ctx context.Context) (int, error) {
	var count int

	err := s.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM albums`,
	).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func nullableTime(t time.Time) any {
	if t.IsZero() {
		return nil
	}

	return t
}
