package poller

import (
	"context"
	"fmt"
	"log"
	"time"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"example.com/subsonic_bot/internal/storage"
	"example.com/subsonic_bot/internal/subsonic"
)

func Run(
	ctx context.Context,
	store *storage.Store,
	subsonicClient *subsonic.Client,
	bot *tgbotapi.BotAPI,
	pollInterval time.Duration,
) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			albums, err := subsonicClient.GetNewestAlbums(ctx, 20)
			if err != nil {
				return err
			}

			subscribers, err := store.ListSubscribers(ctx)
			if err != nil {
				return err
			}

			for _, album := range albums {
				seen, err := store.AlbumIsSeen(ctx, album.ID)
				if err != nil {
					return err
				}
				if seen {
					continue
				}

				text := fmt.Sprintf(
					"На сервер добавлен новый релиз: %s — %s",
					album.Artist,
					album.Title,
				)

				for _, chatID := range subscribers {
					if _, err := bot.Send(tgbotapi.NewMessage(chatID, text)); err != nil {
						log.Printf("there was an error while sending a notification to %d: %v",
								   chatID,
								   err,
						)
					}
				}
				if err := store.SaveSeenAlbum(ctx, album); err != nil {
					return err
				}
			}
		}
	}
}

