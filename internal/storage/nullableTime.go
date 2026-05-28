package storage

import "time"

func nullableTime(t time.Time) any {
	if t.IsZero() {
		return nil
	}

	return t
}

