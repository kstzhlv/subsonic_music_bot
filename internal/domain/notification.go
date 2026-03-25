package domain

import (
	"context"
	"time"
)

type NotificationType string

const (
	NewAlbumNotification NotificationType = "new album"
	AlbumAnniversaryNotification NotificationType = "album anniversary"
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
