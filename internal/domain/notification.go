package domain

import (
	"context"
	"time"
)

type NotificationType string
type NotificationMode string

const (
	NewAlbumNotification NotificationType = "new album"
	AlbumAnniversaryNotification NotificationType = "album anniversary"
)

const (
	NotificationModeAllAlbums NotificationMode = "all"
	NotificationModeWishlistOnly NotificationMode = "wishlist_only"
)

type Notification struct {
	ID int64
	UserID int64
	AlbumID string
	Type NotificationType
	SentAt time.Time
}

type NotificationRepository interface {
	Sent(ctx context.Context, userID int64, albumID string, t NotificationType) (bool, error)
	Save(ctx context.Context, n Notification) error
}
