package util

import (
	"time"

	"github.com/charmbracelet/log"
)

func TimestampToString(timestamp int64) string {
	t, err := time.Unix(timestamp, 0).MarshalText()
	if err != nil {
		log.Error("invalid timestamp", timestamp)
		return "[invalid time]"
	}

	return string(t)
}

func RFC3339ToUnixTimestamp(rfc3339time string) (int64, error) {
	t, err := time.Parse(time.RFC3339, rfc3339time)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}
