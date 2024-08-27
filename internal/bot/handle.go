package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(update tgbotapi.Update) {
	messageText := update.Message.Text
	command := strings.Split(messageText, " ")[0]

	err := b.server.AddUser(update.Message.From.ID)
	if err != nil {
		log.Fatal("failed to add user", "update", update)
	}

	switch command {
	case subscribeCommand:
		b.handleSubscribe(update)
	case unsubscribeCommand:
		b.handleUnsubscribe(update)
	case listCommand:
		b.handleSubscriptionList(update)
	default:
		b.handleHelp(update)
	}
}
