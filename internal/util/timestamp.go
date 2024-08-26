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
