package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"example.com/subsonic_bot/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleNonCommandMessage(
	states map[int64]SessionState,
	update tgbotapi.Update,
	store *storage.Store,
	chatID int64,
	bot *tgbotapi.BotAPI,
) error {
	state := states[chatID]

	switch state.WishlistState {
	case StateWaitingWishlistAdd:
		albumName := strings.TrimSpace(
			update.Message.Text,
		)
		if albumName == "" {
			return nil
		}

		count, err := store.CountWishlistItems(
			context.Background(),
			chatID,
		)
		if err != nil {
			replyWithError(
				bot,
				"count wishlist items",
				err,
				chatID,
				"Произошла ошибка. Попробуйте позже",
			)

			return err
		}

		if count >= 10 {
			_, _ = bot.Send(tgbotapi.NewMessage(
				chatID,
				"В вашем списке желаемого максимальное количество альбомов." +
				" Удалите альбомы из списка желаемого или подождите, пока администраторы добавят релизы из вашего списка.",
			))					
			
			states[chatID] = SessionState{
				WishlistState: StateIdle,
			}

			return nil
		}

		inserted, err := store.AddToWishlist(
			context.Background(),
			chatID,
			albumName,
		)
		if err != nil {
			replyWithError(
				bot,
				"add to wishlist",
				err,
				chatID,
				"Ошибка при добавлении альбома в список желаемого",
			)

			return err
		}

		if !inserted {
			_, _ = bot.Send(tgbotapi.NewMessage(
				chatID,
				fmt.Sprintf(
					"%s уже есть в Вашем списке желаемого!",
					albumName,
				),
			))

			return nil
		}

		_, _ = bot.Send(tgbotapi.NewMessage(
			chatID,
			fmt.Sprintf(
				"%s был добавлен в Ваш список желаемого! " +
				"Админы в скором времени добавят этот релиз!",
				albumName,
			),
		))


		adminChatIDs, err := ParseAdminChatIDs(
			os.Getenv("ADMIN_CHAT_IDS"),
		)
		if err != nil {
			log.Fatal(err)
		}

		for _, id := range adminChatIDs {
			_, _ = bot.Send(tgbotapi.NewMessage(
				id,
				fmt.Sprintf(
					"Пользователем был добавлен альбом %s в вишлист",
					albumName,	
				),
			))
		}

		states[chatID] = SessionState{
			WishlistState: StateIdle,
		}
	}

	return nil
}
