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
		response.Header.Get("x-ratelimit-limit"))
	if err != nil {
		return responseRateLimit{}, err
	}

	remaining, err := strconv.Atoi(
		response.Header.Get("x-ratelimit-remaining"))
	if err != nil {
		return responseRateLimit{}, err
	}

	reset, err := strconv.ParseInt(
		response.Header.Get("x-ratelimit-reset"), 10, 64)
	if err != nil {
		return responseRateLimit{}, err
	}

	resetAfter, err := strconv.ParseFloat(
		response.Header.Get("x-ratelimit-reset-after"), 64)
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

func (c *Client) updateRateLimiter(rl responseRateLimit) {
	c.limiter.SetLimit(rl.limit)
	c.limiter.SetRemaining(rl.remaining - 1)

	// HACK: because i was getting reset time being in the past way too often
	if time.Until(time.Unix(rl.reset, 0)).Nanoseconds() < 0 {
		rl.reset = time.Now().Add(100 * time.Second).Unix()
	}

	c.limiter.SetResetTime(int64(float64(rl.reset)))

	// the client is getting this header only when low on remaining requests,
	// so it might as well don't do any for the time being
	if rl.resetAfter > 0 {
		c.limiter.SetRemaining(0)
		c.limiter.SetResetTime(rl.reset + int64(math.Ceil(rl.resetAfter)))
	}
}
