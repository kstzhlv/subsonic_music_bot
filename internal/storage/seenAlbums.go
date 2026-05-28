package storage

import (
	"context"
	"database/sql"
	"time"

	"example.com/subsonic_bot/internal/domain"
)

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

