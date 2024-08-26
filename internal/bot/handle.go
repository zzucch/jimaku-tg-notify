package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleMessage(update tgbotapi.Update) {
	messageText := update.Message.Text
	command := strings.Split(messageText, " ")[0]

	switch command {
	case subscribeCommand:
		handleSubscribe(update)
	case unsubscribeCommand:
		handleUnsubscribe(update)
	case listCommand:
		handleSubscriptionList(update)
	default:
		handleHelp(update)
	}
}
