package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

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

	if err := store.Init(context.Background()); err != nil {
		log.Fatal(err)
	}

	telegram.ConfigureTelegramBot(
		telegramToken,
		subsonicClient,
		store
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

