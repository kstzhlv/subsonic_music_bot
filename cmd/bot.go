package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"example.com/subsonic_bot/internal/domain"
	"example.com/subsonic_bot/internal/subsonic"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	baseURL := os.Getenv("SUBSONIC_BASE_URL")
	username := os.Getenv("SUBSONIC_USERNAME")
	password := os.Getenv("SUBSONIC_PASSWORD")
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	caCertFile := os.Getenv("SUBSONIC_CA_CERT_FILE")

	if baseURL == "" || username == "" || password == "" {
		log.Fatal("missing SUBSONIC_BASE_URL, SUBSONIC_USERNAME, or SUBSONIC_PASSWORD")
	}

	if telegramToken == "" {
		log.Fatal("missing Telegram bot token")
	}

	httpClient, err := newHTTPClient(caCertFile)
	if err != nil {
		log.Fatal(err)
	}

	subsonicClient := &subsonic.Client{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		HTTPClient: httpClient,
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
		
			_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, formatAlbums(albums)))
		}
	}
}

func newHTTPClient(caCertFile string) (*http.Client, error) {
	if caCertFile == "" {
		return http.DefaultClient, nil
	}

	caCertPEM, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("read SUBSONIC_CA_CERT_FILE: %w", err)
	}

	roots, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("load system cert pool: %w", err)
	}
	if roots == nil {
		roots = x509.NewCertPool()
	}

	if ok := roots.AppendCertsFromPEM(caCertPEM); !ok {
		return nil, fmt.Errorf("parse CA certificate from %s", caCertFile)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		RootCAs: roots,
	}

	return &http.Client{Transport: transport}, nil
}

func formatAlbums(albums []domain.Album) string {
	if len(albums) == 0 {
		return "Альбомов не найдено"
	}

	var b strings.Builder
	for i, album := range albums {
		fmt.Fprintf(&b, "%d. %s - %s\n", i + 1, album.Artist, album.Title)
	}

	return b.String()
}
