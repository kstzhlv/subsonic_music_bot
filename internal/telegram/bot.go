package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"example.com/subsonic_bot/internal/domain"
	"example.com/subsonic_bot/internal/storage"
	"example.com/subsonic_bot/internal/subsonic"
)

func Run(
	bot *tgbotapi.BotAPI,
	subsonicClient *subsonic.Client,
	store *storage.Store,
) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	states := map[int64]SessionState{}

	for update := range updates {
		if update.CallbackQuery != nil {
			err := processCallbackQuery(
				update,
				bot,
				store,
				states,
			)
			
			log.Printf(
				"There was a problem with callback query: %s",
				err,
			)

			continue
		}

		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		state := states[chatID]

		if !update.Message.IsCommand() {
			switch state.WishlistState {
			case StateWaitingWishlistAdd:
				albumName := strings.TrimSpace(
					update.Message.Text,
				)
				if albumName == "" {
					continue
				}

				count, err := store.CountWishlistItems(
					context.Background(),
					chatID,
				)

				if count == 10 {
					_, _ = bot.Send(tgbotapi.NewMessage(
						chatID,
						"В вашем списке желаемого максимальное количество альбомов." +
						" Удалите альбомы из списка желаемого или подождите, пока администраторы добавят релизы из вашего списка.",
					))					

					continue
				}
				if err != nil {
					replyWithError(
						bot,
						"count wishlist items",
						err,
						chatID,
						"Произошла ошибка. Попробуйте позже",
					)

					continue
				}

				err = store.AddToWishlist(
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
					continue
				}

				_, _ = bot.Send(tgbotapi.NewMessage(
					chatID,
					fmt.Sprintf(
						"%s был добавлен в Ваш список желаемого! " +
						"Админы в скором времени добавят этот релиз!",
						albumName,
					),
				))

                states[chatID] = SessionState{
					WishlistState: StateIdle,
				}
			}

			continue
		}

		switch update.Message.Command() {
		case "start":
			err := store.AddSubscriber(
				context.Background(),
				update.Message.Chat.ID,
			)		

			if err != nil {
				replyWithError(
					bot,
					"start subscription",
					err,
					update.Message.Chat.ID,
					"Ошибка при попытке подписаться на рассылку. Попробуйте позже.",
				)
				continue
			}

            continue

		case "stop":
			err := store.RemoveSubscriber(
				context.Background(),
				update.Message.Chat.ID,
			)		

			if err != nil {
				replyWithError(
					bot,
					"stop subscription",
					err,
					update.Message.Chat.ID,
					"Ошибка при попытке отписаться от рассылки. Попробуйте позже.",
				)
				continue
			}

			continue

		case "latest":
			albums, err := subsonicClient.GetNewestAlbums(context.Background(), 20)
			if err != nil {
				replyWithError(
					bot,
					"get latest albums",
					err,
					update.Message.Chat.ID,
					"Ошибка при получении списка альбомов.",
				)
				continue
			}

			_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, domain.FormatAlbums(albums)))

			continue

		case "id":
			text := fmt.Sprintf(
				"chat_id=%d\nfrom_id=%d",
				update.Message.Chat.ID,
				update.Message.From.ID,
			)
			_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))

			continue

		case "wishlist":
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
					"Вывести список желаемого",
					"wishlist_show"),
				),

				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
					"Добавить в список желаемого",
					"wishlist_add"),
				),

				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
					"Удалить из списка желаемого",
					"wishlist_remove"),
				),	

				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
					"Очисить список желаемого",
					"wishlist_empty"),
				),	
			)
			
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Что Вы хотите сделать?",
			)
			msg.ReplyMarkup = keyboard	

			_, err := bot.Send(msg)
			if err != nil {
				replyWithError(
					bot,
					"wishlist",
					err,
					update.Message.Chat.ID,
					"Ошибка при работе со списком желаемого.",
				)
				continue
			}

			continue
		}
	}
}

