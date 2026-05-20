package domain

import "time"

type WishlistItem struct {
	ID int64
	ChatID int64
	AlbumName string
	CreatedAt time.Time
}
