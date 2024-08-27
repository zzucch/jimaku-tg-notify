package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/zzucch/jimaku-tg-notify/internal/util"
	"golang.org/x/time/rate"
)

const (
	rateLimitBurst    = 25
	rateLimitInterval = 100 * time.Second
)

type Entry struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	LastModified string `json:"last_modified"`
	AnilistID    int64  `json:"anilist_id"`
	EnglishName  string `json:"english_name"`
	JapaneseName string `json:"japanese_name"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
	limiter    *rate.Limiter
	mu         sync.Mutex
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		limiter: rate.NewLimiter(
			rate.Every(
				rateLimitInterval/time.Duration(rateLimitBurst)),
			rateLimitBurst),
	}
}

func (c *Client) GetLatestSubtitle(titleID int64) (int64, error) {
	entryData, err := c.GetEntryData(titleID)
	if err != nil {
		return 0, err
	}

	latestSubtitleTime, err := util.RFC3339ToUnixTimestamp(
		entryData.LastModified)
	if err != nil {
		return 0, err
	}

	return latestSubtitleTime, nil
}

func (c *Client) GetEntryData(titleID int64) (*Entry, error) {
	url := "https://jimaku.cc/api/entries/" +
		strconv.FormatInt(titleID, 10)

	response, err := c.getResponse(url)
	if err != nil {
		return nil, err
	}

	var entry Entry
	if err = json.Unmarshal([]byte(response), &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

func (c *Client) getResponse(url string) (string, error) {
	if err := c.limiter.Wait(context.Background()); err != nil {
		return "", err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("Authorization", c.apiKey)
	request.Header.Add("Accept", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New(
			"unexpected status code: " + strconv.Itoa(response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
