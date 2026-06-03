package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
			err := ProcessCallbackQuery(
				update,
				bot,
				store,
				states,
			)
			
			if err != nil {
				log.Printf(
					"There was a problem with callback query: %s",
					err,
				)
			}

			continue
		}

		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		if !update.Message.IsCommand() {
			err := HandleNonCommandMessage(
				states,
				update,
				store,
				chatID,
				bot,
				)

			if err != nil {
				log.Printf(
					"There was a problem with non command message: %s",
					err,
				)
			}
			continue
		}

		err := HandleCommandMessage(
			update,
			store,
			bot,
			subsonicClient,
		)

		if err != nil {
			log.Printf(
				"There was a problem with command message: %s",
				err,
			)
		}
	}
}

