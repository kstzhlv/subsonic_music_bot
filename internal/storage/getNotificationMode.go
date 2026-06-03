package storage

import (
	"context"
	"database/sql"

	"example.com/subsonic_bot/internal/domain"
)

func (s *Store) GetNotificationMode(
	ctx context.Context,
	chatID int64,
) (domain.NotificationMode, error) {
	var mode string	

	err := s.db.QueryRowContext(
		ctx,
		`
		SELECT mode
		FROM notification_settings
		WHERE chat_id = ?
		`,
		chatID,
	).Scan(&mode)

	if err == sql.ErrNoRows {
		return domain.NotificationModeAllAlbums, nil
	}

	if err != nil {
		return "", err
	}

	return domain.NotificationMode(mode), nil
}

func (s *Store) 
