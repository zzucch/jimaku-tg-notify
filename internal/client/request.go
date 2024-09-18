package client

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
)

func (c *Client) getResponse(url string, attempts int) (string, error) {
	for attempt := range attempts {
		c.limiter.Wait()

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}

		request.Header.Add("Authorization", c.apiKey)
		request.Header.Add("Accept", "application/json")

		response, err := c.httpClient.Do(request)
		if err != nil {
			if attempt < attempts-1 {
				log.Warn(
					"failed to do http request",
					"url",
					url,
					"attempt",
					attempt+1,
					"err",
					err)
				continue
			}

			return "", err
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusTooManyRequests {
			log.Warn(
				"too many requests",
				"url",
				url,
				"attempt",
				attempt+1,
			)

			rateLimit, err := c.parseRateLimitHeaders(response)
			if err != nil {
				return "", err
			}

			c.updateRateLimiter(rateLimit)

			continue
		}

		if response.StatusCode == http.StatusOK {
			rateLimit, err := c.parseRateLimitHeaders(response)
			if err != nil {
				return "", err
			}

			c.updateRateLimiter(rateLimit)

			body, err := io.ReadAll(response.Body)
			if err != nil {
				return "", err
			}

			return string(body), nil
		}

		message := "Unexpected response status code: " +
			strconv.Itoa(response.StatusCode)

		if response.StatusCode == http.StatusUnauthorized {
			message = "Consider checking your API key"
		}

		if response.StatusCode == http.StatusNotFound {
			message = "The entry does not exist"
		}

		return "", errors.New(message)
	}

	return "", errors.New("exceeded maximum retry attempts")
}
