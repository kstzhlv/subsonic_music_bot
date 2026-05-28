package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

