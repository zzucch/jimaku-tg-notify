package timeutil

import (
	"time"

	"github.com/charmbracelet/log"
)

func TimestampToString(timestamp int64) string {
	utcTime := time.Unix(timestamp, 0).UTC()
	formattedTime := utcTime.Format(time.DateTime)

	if len(formattedTime) != len(time.DateTime) {
		log.Error(
			"cannot convert timestamp to string",
			"timestamp",
			timestamp,
			"formattedTime",
			formattedTime,
		)

		return "[invalid time]"
	}

	return formattedTime
}

func RFC3339ToUnixTimestamp(rfc3339time string) (int64, error) {
	timestamp, err := time.Parse(time.RFC3339, rfc3339time)
	if err != nil {
		return 0, err
	}

	return timestamp.Unix(), nil
}

func AddUTCOffsetInMinutes(inputTime time.Time, offsetMinutes int) time.Time {
	duration := time.Duration(offsetMinutes) * time.Minute
	return inputTime.Add(duration)
}
