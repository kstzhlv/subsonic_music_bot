package storage

import (
	"context"
	"time"
)

func (s *Store) AddToWishlist(
	ctx context.Context,
	chatID int64,
	albumName string,
) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT OR IGNORE INTO wishlist_items
		(chat_id, album_name, created_at)
		VALUES (?, ?, ?)`,
		chatID,
		albumName,
		time.Now().UTC(),
	)

	return err
}

func (s *Store) RemoveFromWishlist(
	ctx context.Context,
	chatID int64,
	itemID int64,
) error {
	_, err := s.db.ExecContext(
		ctx,
		`DELETE FROM wishlist_items
		WHERE chat_id = ?
		AND id = ?
		`,
		chatID,
		itemID,
	)

	return err
}

func (s *Store) ClearWishlist(
	ctx context.Context,
	chatID int64,
) error {
	_, err := s.db.ExecContext(
		ctx,
		`
		DELETE FROM wishlist_items
		WHERE chat_id = ?
		`,
		chatID,	
	)

	return err
}

