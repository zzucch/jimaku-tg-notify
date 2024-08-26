package client

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/parse"
)

func GetLatestSubtitleTimestamp(titleID int64) (int64, error) {
	dates, err := getDates(
		"https://jimaku.cc/entry/" + strconv.FormatInt(titleID, 10))
	if err != nil {
		log.Error("failed to get subtitle dates", "titleID", titleID, "err", err)
    return -1, err
	}

	var latest int64 = -1
	for _, date := range dates {
		timestamp := date.ToUnixTimestamp()
		if timestamp > latest {
			latest = timestamp
		}
	}

	return latest, nil
}

func getDates(url string) ([]parse.Date, error) {
	response, err := getResponse(url)
	if err != nil {
		return nil, err
	}

	dates := parse.GetDates(strings.Split(response, "\n"))
	return dates, nil
}

func getResponse(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("unexpected status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
