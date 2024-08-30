package bot

import (
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const maxRetries = 5

func (b *Bot) SendMessage(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)

	for retry := 0; retry < maxRetries; retry++ {
		if _, err := b.botAPI.Send(message); err != nil {
			if retryAfter, ok := getRetryAfterDuration(err); ok {
				retryAfter *= 2
				log.Warn("too many requests to telegram", "retry after", retryAfter)
				time.Sleep(retryAfter)

				continue
			}

			log.Error("failed to send message", "err", err, "retry", retry)
		}

		return
	}

	log.Error("failed to send message", "max retries", maxRetries)
}

func getRetryAfterDuration(err error) (time.Duration, bool) {
	if err == nil {
		return 0, false
	}

	const prefix = "Too Many Requests: retry after "

	errMsg := err.Error()
	if strings.Contains(errMsg, prefix) {
		parts := strings.Split(errMsg, prefix)

		const correctPartsAmount = 2

		if len(parts) == correctPartsAmount {
			if retryAfterSeconds, parseErr :=
				strconv.Atoi(parts[1]); parseErr == nil {
				return time.Duration(retryAfterSeconds) * time.Second, true
			}
		}
	}

	return 0, false
}
