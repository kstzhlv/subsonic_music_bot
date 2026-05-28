package telegram

import (
	"fmt"
	"strings"

	"example.com/subsonic_bot/internal/domain"
)

func formatWishlistItems(
	items []domain.WishlistItem,
) string {
	if len(items) == 0 {
		return "Список желаемого пуст"
	}	

	var b strings.Builder
	b.WriteString("Ваш список желаемого:\n")

	for i, item := range items {
		fmt.Fprintf(
			&b,
			"%d. %s\n",
			i + 1,
			item.AlbumName,
		)
	}

	return b.String()
}

