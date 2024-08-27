package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleNotificationIntervalChange(update tgbotapi.Update) {
	chatID := update.Message.From.ID

	unvalidatedTitleID := update.Message.Text[len(intervalCommand):]
	unvalidatedTitleID = strings.TrimSpace(unvalidatedTitleID)

	interval, err := strconv.Atoi(unvalidatedTitleID)
	if unvalidatedTitleID == "" || err != nil {
		b.SendMessage(chatID, "invalid command")
		return
	}

	if err := b.server.SetInterval(chatID, interval); err != nil {
		b.SendMessage(chatID, "failed to process: "+err.Error())
    return
	}

	b.SendMessage(chatID, "done")
}
