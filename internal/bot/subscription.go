package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleSubscribe(update tgbotapi.Update) {
	b.handleSubscription(update, subscribeCommand, b.server.Subscribe)
}

func (b *Bot) handleUnsubscribe(update tgbotapi.Update) {
	b.handleSubscription(update, unsubscribeCommand, b.server.Unsubscribe)
}

func (b *Bot) handleSubscription(
	update tgbotapi.Update,
	command string,
	action func(chatID int64, titleID int64) error,
) {
	chatID := update.Message.Chat.ID

	unvalidatedTitleID := update.Message.Text[len(command):]
	unvalidatedTitleID = strings.TrimSpace(unvalidatedTitleID)

	titleID, err := strconv.ParseInt(unvalidatedTitleID, 10, 64)
	if unvalidatedTitleID == "" || err != nil {
		b.SendMessage(chatID, "invalid command")
		return
	}

	if err := action(chatID, titleID); err != nil {
		b.SendMessage(
			update.Message.From.ID,
			"failed to process: "+err.Error())

		return
	}

	b.SendMessage(update.Message.From.ID, "done")
}