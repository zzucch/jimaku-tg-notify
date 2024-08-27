package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(update tgbotapi.Update) {
	messageText := update.Message.Text
	command := strings.Split(messageText, " ")[0]

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
