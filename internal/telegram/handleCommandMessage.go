package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"example.com/subsonic_bot/internal/domain"
	"example.com/subsonic_bot/internal/storage"
	"example.com/subsonic_bot/internal/subsonic"
)

func HandleCommandMessage(
	update tgbotapi.Update,
	store *storage.Store,
	bot *tgbotapi.BotAPI,
	subsonicClient *subsonic.Client,
) error {
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
			return err
		}

		return nil

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
			return err
		}

		return nil

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
			return err
		}

		_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, domain.FormatAlbums(albums)))

		return nil

	case "id":
		text := fmt.Sprintf(
			"chat_id=%d\nfrom_id=%d",
			update.Message.Chat.ID,
			update.Message.From.ID,
		)
		_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))

		return nil

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
			return err
		}

		return nil

	case "notifications":
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
				"Получать уведомления о любых новых альбомах",
				"notifications_all"),
			),

			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
				"Получать уведомления только об альбомах из своего списка желаний",
				"notifications_wishlist"),
			),
		)
		
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Выберите формат уведомлений",
		)
		msg.ReplyMarkup = keyboard	

		_, err := bot.Send(msg)
		if err != nil {
			replyWithError(
				bot,
				"notifications configuration",
				err,
				update.Message.Chat.ID,
				"Ошибка при работе с уведомлениями.",
			)
			return err
		}

	}

	return nil
}
