package storage

import (
	"context"
	"database/sql"

	"example.com/subsonic_bot/internal/domain"
)

func (s *Store) ListWishlistItems(
	ctx context.Context,
	chatID int64,
) ([]domain.WishlistItem, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT id, chat_id, album_name, created_at
		FROM wishlist_items
		WHERE chat_id = ?
		ORDER BY created_at
		`,
		chatID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlistItems []domain.WishlistItem
	for rows.Next() {
		var item domain.WishlistItem
		if err := rows.Scan(
			&item.ID,
			&item.ChatID,
			&item.AlbumName,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		wishlistItems = append(wishlistItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wishlistItems, nil
}

func (s *Store) CountWishlistItems(
	ctx context.Context,
	chatID int64,
) (int, error) {
	var count int

	err := s.db.QueryRowContext(
		ctx,
		`
		SELECT COUNT(*)
		FROM wishlist_items
		WHERE chat_id = ?
		`,
		chatID,
	).Scan(&count)

	if err == sql.ErrNoRows {
		return 0, err
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}
