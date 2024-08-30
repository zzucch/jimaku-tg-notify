package client

import (
	"math"
	"net/http"
	"strconv"
	"time"
)

type responseRateLimit struct {
	limit      int
	remaining  int
	reset      int64
	resetAfter float64
}

func (c *Client) parseRateLimitHeaders(
	response *http.Response,
) (responseRateLimit, error) {
	limit, err := strconv.Atoi(
		response.Header.Get("X-Ratelimit-Limit"))
	if err != nil {
		return responseRateLimit{}, err
	}

	remaining, err := strconv.Atoi(
		response.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return responseRateLimit{}, err
	}

	reset, err := strconv.ParseInt(
		response.Header.Get("X-Ratelimit-Reset"), 10, 64)
	if err != nil {
		return responseRateLimit{}, err
	}

	resetAfter, err := strconv.ParseFloat(
		response.Header.Get("X-Ratelimit-Reset-After"), 64)
	if err != nil {
		resetAfter = 0
	}

	return responseRateLimit{
		limit:      limit,
		remaining:  remaining,
		reset:      reset,
		resetAfter: resetAfter,
	}, nil
}

func (c *Client) updateRateLimiter(responseRateLimit responseRateLimit) {
	c.limiter.SetLimit(responseRateLimit.limit)
	c.limiter.SetRemaining(responseRateLimit.remaining - 1)

	// HACK: because i was getting reset time being in the past way too often
	if time.Until(time.Unix(responseRateLimit.reset, 0)).Nanoseconds() < 0 {
		responseRateLimit.reset = time.Now().Add(100 * time.Second).Unix()
	}

	c.limiter.SetResetTime(int64(float64(responseRateLimit.reset)))

	// the client is getting this header only when low on remaining requests,
	// so it might as well don't do any for the time being
	if responseRateLimit.resetAfter > 0 {
		c.limiter.SetResetTime(
			responseRateLimit.reset +
				int64(math.Ceil(responseRateLimit.resetAfter)))

		c.limiter.SetRemaining(0)
	}
}
