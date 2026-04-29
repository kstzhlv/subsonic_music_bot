package telegram

import (
	"context"
	"log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"example.com/subsonic_bot/internal/subsonic"
	"example.com/subsonic_bot/internal/domain"
)

func ConfigureTelegramBot(telegramToken string, subsonicClient *subsonic.Client) {
	if telegramToken == "" {
		log.Fatal("missing Telegram bot token")
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "latest":
			albums, err := subsonicClient.GetNewestAlbums(context.Background(), 18)
			if err != nil {
				log.Printf("latest command failed for chat %d: %v", update.Message.Chat.ID, err)
				_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при загрузке списка альбомов"))
				continue
			}
		
			_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, domain.FormatAlbums(albums)))
		}
	}
}
