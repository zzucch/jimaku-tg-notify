package util

import (
	"time"

	"github.com/charmbracelet/log"
)

func TimestampToString(timestamp int64) string {
	text, err := time.Unix(timestamp, 0).UTC().MarshalText()
	if err != nil {
		log.Error(
			"cannot convert timestamp to string",
			"timestamp",
			timestamp)

		return "[invalid time]"
	}

	return string(text)
}

func RFC3339ToUnixTimestamp(rfc3339time string) (int64, error) {
	timestamp, err := time.Parse(time.RFC3339, rfc3339time)
	if err != nil {
		return 0, err
	}

	return timestamp.Unix(), nil
}
