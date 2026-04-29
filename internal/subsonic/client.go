package subsonic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"example.com/subsonic_bot/internal/domain"
)

type Client struct {
	BaseURL string
	Username string
	Password string
	HTTPClient *http.Client
}

func (c *Client) GetNewestAlbums(ctx context.Context, size int) ([]domain.Album, error) {
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
	if err != nil {
		return nil, err
	}

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("subsonic returned status %d", response.StatusCode)
	}

	var parsed albumListResponse
	if err := json.NewDecoder(response.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	subsonicResp := parsed.SubsonicResponse
	if subsonicResp.Status != "ok" {
		return nil, fmt.Errorf("subsonic error: %s", subsonicResp.Error.Message)
	}

	albums := make([]domain.Album, 0, len(subsonicResp.AlbumList2.Album))
	for _, a := range subsonicResp.AlbumList2.Album {
		albums = append(albums, domain.Album{
			ID: a.ID,
			Title: a.Name,
			Artist: a.Artist,
			AddedAt: parseSubsonicTime(a.Created),
		})
	}

	return albums, nil
}


func subsonicToken(password, salt string) string {
	sum := md5.Sum([]byte(password + salt))
	return hex.EncodeToString(sum[:])
}

func parseSubsonicTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return t
}

type albumListResponse struct {
	SubsonicResponse subsonicResponse `json:"subsonic-response"`
}

type subsonicResponse struct {
	Status string `json:"status"`
	Version string `json:"version"`
	AlbumList2 albumList2 `json:"albumList2"`
	Error apiError `json:"error"`
}

type albumList2 struct {
	Album []apiAlbum `json:"album"`
}

type apiAlbum struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Artist string `json:"artist"`
	Created string `json:"created"`
}

type apiError struct {
	Code int `json:"code"`
	Message string `json:"message"`
}
