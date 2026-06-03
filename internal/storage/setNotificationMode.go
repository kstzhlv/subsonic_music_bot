package storage

import (
	"context"

	"example.com/subsonic_bot/internal/domain"
)

func (s *Store) SetNotificationMode(
	ctx context.Context,
	chatID int64,
	mode domain.NotificationMode,
) error {
	_, err := s.db.ExecContext(
		ctx,
		`
		INSERT INTO notification_settings (chat_id, mode)
		VALUES (?, ?)
		ON CONFLICT(chat_id)
		DO UPDATE SET mode = excluded.mode
		`,
		chatID,
		string(mode),
	)

	return err
}
