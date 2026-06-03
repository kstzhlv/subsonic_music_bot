package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
		
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"example.com/subsonic_bot/internal/storage"
)

func ProcessCallbackQuery(
	update tgbotapi.Update,
	bot *tgbotapi.BotAPI,
	store *storage.Store,
	states map[int64]SessionState,
) error {
	_, _ = bot.Request(tgbotapi.NewCallback(
		update.CallbackQuery.ID,
		"",
	))

	chatID := update.CallbackQuery.Message.Chat.ID
	data := update.CallbackQuery.Data

	switch {
	case data == "wishlist_show":
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
			return err
		}

		_, _ = bot.Send(tgbotapi.NewMessage(
			chatID,
			formatWishlistItems(items),	
		))

		return nil

	case data == "wishlist_add":
		states[chatID] = SessionState{
			WishlistState: StateWaitingWishlistAdd,
		}

		_, _ = bot.Send(tgbotapi.NewMessage(
			chatID,
			"Введите название альбома, который Вы бы хотели видеть на сервере",
		))

		return nil

	case data == "wishlist_remove":
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

			return err
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
			"Выберите альбом, который хотите удалить из Вашего списка желаемого",
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

			return err
		}

		return nil


	case data == "wishlist_empty":
		err := store.ClearWishlist(
			context.Background(),
			chatID,
		)	
		if err != nil {
			replyWithError(
				bot, 
				"wishlist_empty",
				err,
				chatID,
				"Ошибка при очистке списка желаемого",
			)

			return err
		}

		_, _ = bot.Send(tgbotapi.NewMessage(
			chatID,
			"Ваш список желаемого очищен",
		))

		return nil

	case strings.HasPrefix(
		data,
		"wishlist_remove:",
	):
		idStr := strings.TrimPrefix(
			data,
			"wishlist_remove:",
		)	
		itemID, err := strconv.ParseInt(
			idStr,
			10,
			64,
		)
		if err != nil {
			replyWithError(
				bot,
				"parse wishlist item ID",
				err,
				chatID,
				"Некорректный идентификатор альбома.",
			)

			return err
		}

		err = store.RemoveFromWishlist(
			context.Background(),
			chatID,
			itemID,
		)

		if err != nil {
			replyWithError(
				bot,
				"remove item from wishlist",
				err,
				chatID,
				"Ошибка при удалении из списка желаемого.",
			)

			return err
		}

		_, _ = bot.Send(tgbotapi.NewMessage(
			chatID,
			"Альбом удалён из списка желаемого",
		))

		return nil

	case strings.HasPrefix(data, "notifications"):
		if data == "notifications_all" {

		}
	}

	return nil
}

