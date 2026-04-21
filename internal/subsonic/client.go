package subsonic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"

	"example.com/subsonic_bot/internal/domain"
)

type Client struct {
	BaseURL string
	Username string
	Password string
	HTTPClient *http.client
}

func (c *Client) getNewestAlbums(ctx context.Context, size int)
([]domain.Album, error) {
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	salt := "subsonictelegrambot"
	token := subsonicToken(c.Password, salt)

	endpoint, err := url.Parse(c.BaseURL + "/rest/getAlbumList2.view")
	if err != nil {
		return nil, err
	}

	q := endpoint.Query()
	q.Set("u", c.Username)
	q.Set("t", token)
	q.Set("s", salt)
	q.Set("v", "1.16.1")
	q.Set("c", "subsonic_bot")
	q.Set("f", "json")
	q.Set("type", "newest")
	q.Set("size", fmt.Sprintf("%d", size))
	endpoint.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
}

func subsonicToken(password, salt string) string {
	sum := md5.Sum([]byte(password + salt))
	return hex.EncodeToString(sum[:])
}
