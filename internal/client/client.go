package client

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zzucch/jimaku-tg-notify/internal/dto"
	"github.com/zzucch/jimaku-tg-notify/internal/rate"
)

const attemptsAmount = 5

type Client struct {
	apiKey     string
	httpClient *http.Client
	limiter    *rate.Limiter
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		limiter:    rate.NewLimiter(),
	}
}

func (c *Client) UpdateAPIKey(apiKey string) {
	c.apiKey = apiKey
}

func (c *Client) GetEntryDetails(titleID int64) (*dto.Entry, error) {
	url := "https://jimaku.cc/api/entries/" +
		strconv.FormatInt(titleID, 10)

	response, err := c.getResponse(url, attemptsAmount)
	if err != nil {
		return nil, err
	}

	var entry dto.Entry
	if err = json.Unmarshal([]byte(response), &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}
