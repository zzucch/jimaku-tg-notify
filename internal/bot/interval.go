package bot

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleNotificationIntervalChange(update tgbotapi.Update) {
	chatID := update.Message.From.ID

	unvalidatedInterval := update.Message.Text[len(intervalCommand):]
	unvalidatedInterval = strings.TrimSpace(unvalidatedInterval)

	interval, err := strconv.Atoi(unvalidatedInterval)
	if unvalidatedInterval == "" || err != nil {
		if err := b.SendMessage(
			chatID,
			"Example usage:\n"+intervalCommand+" 24",
		); err != nil {
			log.Error("failed to send message", "err", err)
		}
		return
	}

	const maxSupportedInterval = int(math.MaxInt64 / time.Hour)
	if interval > maxSupportedInterval {
		if err := b.SendMessage(
			chatID,
			"Do not use unreasonably long interval",
		); err != nil {
			log.Error("failed to send message", "err", err)
		}
		return
	}

	if err := b.server.SetInterval(chatID, interval); err != nil {
		if err := b.SendMessage(
			chatID,
			"Failed to process.\n"+err.Error(),
		); err != nil {
			log.Error("failed to send message", "err", err)
		}
		return
	}

	if err := b.SendMessage(
		chatID,
		"Notification interval is set",
	); err != nil {
		log.Error("failed to send message", "err", err)
	}
}
