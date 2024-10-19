package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleSettingUTCOffset(update tgbotapi.Update) {
	chatID := update.Message.From.ID

	unvalidatedOffset := update.Message.Text[len(offsetCommand):]
	unvalidatedOffset = strings.TrimSpace(unvalidatedOffset)

	offset, err := strconv.Atoi(unvalidatedOffset)
	if unvalidatedOffset == "" || err != nil {
		_ = b.SendMessage(chatID, "Example usage:\n"+offsetCommand+" 330")
		return
	}

	const (
		maxSupportedOffset = 14 * 60
		minSupportedOffset = -12 * 60
	)

	if minSupportedOffset > offset || offset > maxSupportedOffset {
		_ = b.SendMessage(
			chatID,
			"Failed to process. The UTC offset in minutes must be between -720 and 840",
		)

		return
	}

	if err := b.server.SetUTCOffset(chatID, offset); err != nil {
		_ = b.SendMessage(chatID, "Failed to process.\n"+err.Error())
		return
	}

	_ = b.SendMessage(chatID, "UTC offset in minutes is set")
}
