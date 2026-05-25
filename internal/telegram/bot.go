package telegram

import (
	"context"
	"fmt"
	"log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"example.com/subsonic_bot/internal/storage"
	"example.com/subsonic_bot/internal/subsonic"
	"example.com/subsonic_bot/internal/domain"
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
			_, _ = bot.Request(tgbotapi.NewCallback(
				update.CallbackQuery.ID,
				"",
			))

			chatID := update.CallbackQuery.Message.Chat.ID

			switch update.CallbackQuery.Data {
			case "wishlist_show":
				items, err := store.ListWishlistItems(
					context.Background(),
					chatID,
				)
				if err != nil {
					replyWithError(
						bot,
						"Вывод альбомов из списка желаемого",
						err,
						chatID,
						"Ошибка при получении списка альбомов",
					)
				}

				_, _ = bot.Send(tgbotapi.NewMessage(
					chatID,
					"Ваш список желаемого:",
				))

			case "wishlist_add":
				states[chatID] = SessionState{
					WishlistState: StateWaitingWishlistAdd,
				}

				_, _ = bot.Send(tgbotapi.NewMessage(
					chatID,
					"Введите название альбома, который Вы бы хотели видеть на сервере",
				))

				continue

			case "wishlist_remove":
				items, err := store.ListWishlistItems(
					context.Background(),
					chatID,
				)
				if err != nil {
					replyWithError(
						bot,
						"Вывод альбомов из списка желаемого при wishlist_remove",
						err,
						chatID,
						"Ошибка при получении списка альбомов",
					)
				}
				
				keyboard := tgbotapi.NewInlineKeyboardMarkup()
				for _, item := range items {
					row := tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(
							item.AlbumName,
							fmt.Sprintf("wishlist_remove:%d", item.ID),
						),
					)

					keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
				}

				msg := tgbotapi.NewMessage(
					chatID,
					"Выберите номер альбома, который хотите удалить из Вашего списка желаемого",
				)
				msg.ReplyMarkup = keyboard
				_, err = bot.Send(msg)
				if err != nil {
					replyWithError(
						bot, 
						"wishlist_remove send keyboard",
						err,
						chatID,
						"Ошибка при удалении альбома из списка желаемого",
					)
				}


			case "wishlist_empty":
				_, _ = bot.Send(tgbotapi.NewMessage(
					chatID,
					"Ваш список желаемого очищен",
				))
			}
		}

		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		state := states[chatID]

		if !update.Message.IsCommand() {
			switch state.WishlistState {
			case StateWaitingWishlistAdd:
				albumName := update.Message.Text
				err := store.AddToWishlist(
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
				}

				_, _ = bot.Send(tgbotapi.NewMessage(
					chatID,
					"%s был добавлен в Ваш список желаемого! Админы в скором времени добавят этот релиз!",
				))
			}

			states[chatID] = SessionState{
				WishlistState: StateIdle,
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

		case "id":
			text := fmt.Sprintf(
				"chat_id=%d\nfrom_id=%d",
				update.Message.Chat.ID,
				update.Message.From.ID,
			)
			_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))

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
		}
	}
}

func replyWithError(
	bot *tgbotapi.BotAPI,
	operation string,
	err error,
	chatID int64,
	message string,
) {
	log.Printf(
		"%s command failed for chat %d: %v",
		operation,
		chatID,
		err,
	)
	if _, sendError := bot.Send(tgbotapi.NewMessage(
		chatID,
		message,
	)); sendError != nil {
		log.Printf(
			"send error: chat %d: %v",
			chatID,
			sendError,
		)
	}
}
