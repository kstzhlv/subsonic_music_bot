package domain

import (
	"fmt"
	"strings"
)

func FormatAlbums(albums []Album) string {
	if len(albums) == 0 {
		return "Альбомов не найдено"
	}

	var b strings.Builder
	for i, album := range albums {
		fmt.Fprintf(&b, "%d. %s — %s\n", i + 1, album.Artist, album.Title)
	}

	return b.String()
}
