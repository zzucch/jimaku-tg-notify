package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleNotificationIntervalChange(update tgbotapi.Update) {
	chatID := update.Message.From.ID

	unvalidatedInterval := update.Message.Text[len(intervalCommand):]
	unvalidatedInterval = strings.TrimSpace(unvalidatedInterval)

	interval, err := strconv.Atoi(unvalidatedInterval)
	if unvalidatedInterval == "" || err != nil {
		b.SendMessage(chatID, "Example usage:\n"+intervalCommand+" 24")
		return
	}

	if err := b.server.SetInterval(chatID, interval); err != nil {
		b.SendMessage(chatID, "Failed to process.\n"+err.Error())
		return
	}

	b.SendMessage(chatID, "done")
}
