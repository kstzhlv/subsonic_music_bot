package domain

import "context"
import "time"

type Album struct {
	ID string
	Title string
	Artist string
	ReleaseDate time.Time
	AddedAt time.Time
}

type AlbumRepository interface {
	Save(ctx context.Context, album Album) error
	Exists(ctx context.Context, id string) (bool, error)
	GetAlbumByDate(ctx context.Context, date time.Time) ([]Album, error)
}
