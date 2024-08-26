package parse

import "time"

type Date struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

func (d Date) ToUnixTimestamp() int64 {
	return time.Date(
		d.Year,
		time.Month(d.Month),
		d.Day,
		d.Hour,
		d.Minute,
		d.Second,
		0,
		time.UTC).Unix()
}

func GetDates(data []string) []Date {
	dates := make([]Date, 0)
	for _, line := range data {
		tokens, err := lexer(line)
		if err != nil {
			continue
		}

		date, err := parser(tokens)
		if err != nil {
			continue
		}

		dates = append(dates, date)
	}

	return dates
}
