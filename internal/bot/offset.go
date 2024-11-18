package bot

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleSettingUTCOffset(update tgbotapi.Update) {
	chatID := update.Message.From.ID

	unvalidatedOffset := update.Message.Text[len(offsetCommand):]
	unvalidatedOffset = strings.TrimSpace(unvalidatedOffset)

	offset, err := strconv.Atoi(unvalidatedOffset)
	if unvalidatedOffset == "" || err != nil {
		if err := b.SendMessage(
			chatID,
			"Example usage:\n"+offsetCommand+" 330",
		); err != nil {
			log.Error("failed to send message", "err", err)
		}

		return
	}

	const (
		maxSupportedOffset = 14 * 60
		minSupportedOffset = -12 * 60
	)

	if minSupportedOffset > offset || offset > maxSupportedOffset {
		if err := b.SendMessage(
			chatID,
			"Failed to process. The UTC offset in minutes must be between -720 and 840",
		); err != nil {
			log.Error("failed to send message", "err", err)
		}

		return
	}

	if err := b.server.SetUTCOffset(chatID, offset); err != nil {
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
		"UTC offset in minutes is set",
	); err != nil {
		log.Error("failed to send message", "err", err)
	}
}
