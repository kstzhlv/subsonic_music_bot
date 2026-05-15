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

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
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

			_, _ = bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Вы успешно подписались на рассылку о новых альбомах",
			))

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

			_, _ = bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Вы успешно отписались от рассылки о новых альбомах",
			))

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
			text := fmt.Sprintf("chat_id=%d\nfrom_id=%d",
								update.Message.Chat.ID,
								update.Message.From.ID,
								)
			_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))

		case "wishlist":
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
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
