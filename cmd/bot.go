package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "modernc.org/sqlite"

	"example.com/subsonic_bot/internal/poller"
	"example.com/subsonic_bot/internal/storage"
	"example.com/subsonic_bot/internal/subsonic"
	"example.com/subsonic_bot/internal/telegram"
)

func main() {
	baseURL := os.Getenv("SUBSONIC_BASE_URL")
	username := os.Getenv("SUBSONIC_USERNAME")
	password := os.Getenv("SUBSONIC_PASSWORD")
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	caCertFile := os.Getenv("SUBSONIC_CA_CERT_FILE")
	adminChatIDs, err := telegram.ParseAdminChatIDs(
		os.Getenv("ADMIN_CHAT_IDS"),
	)	

	if baseURL == "" || username == "" || password == "" {
		log.Fatal("missing SUBSONIC_BASE_URL, SUBSONIC_USERNAME, or SUBSONIC_PASSWORD")
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

	db, err := sql.Open("sqlite", "/data/subsonic-bot.db")
	if err != nil {
		log.Fatal(err)
	}

	store := storage.New(db)
	ctx := context.Background()

	err = store.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if telegramToken == "" {
		log.Fatal("missing Telegram bot token")
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	count, err := store.CountSeenAlbums(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		albums, err := subsonicClient.GetNewestAlbums(ctx, 20)
		if err != nil {
			log.Fatal(err)
		}
		
		for _, album := range albums {
			if err := store.SaveSeenAlbum(ctx, album);
			err != nil {
				log.Fatal(err)
			}
		}
	}

	go func() {
		if err := poller.Run(
			ctx,
			store,
			subsonicClient,
			bot,
			5 * time.Minute,
		); err != nil {
			log.Fatal(err)
		}
	}()

	telegram.Run(
		bot,
		subsonicClient,
		store,
		adminChatIDs,
	)
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

