package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NotifyAdmins(
	bot *tgbotapi.BotAPI,
	userChatID int64,
	albumName string,
) {
	_, _ = bot.Send(tgbotapi.NewMessage(
		chatID,
		fmt.Sprintf(
			"Пользователь %d добавил в свой вишлист альбом %s",
			userChatID,
            albumName,
		),
	))
}
