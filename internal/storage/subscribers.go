package storage

import (
	"context"
	"time"
)

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

