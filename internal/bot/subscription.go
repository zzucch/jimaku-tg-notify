package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/server"
)

func handleSubscribe(update tgbotapi.Update) {
	handleSubscription(update, subscribeCommand, server.Subscribe)
}

func handleUnsubscribe(update tgbotapi.Update) {
	handleSubscription(update, unsubscribeCommand, server.Unsubscribe)
}

func handleSubscription(
	update tgbotapi.Update,
	command string,
	action func(chatID int64, titleID int64) error,
) {
	chatID := update.Message.Chat.ID

	unvalidatedTitleID := update.Message.Text[len(command):]
	unvalidatedTitleID = strings.TrimSpace(unvalidatedTitleID)

	titleID, err := strconv.ParseInt(unvalidatedTitleID, 10, 64)
	if unvalidatedTitleID == "" || err != nil {
		SendMessage(chatID, "invalid command")
		return
	}

	if err := action(chatID, titleID); err != nil {
		SendMessage(update.Message.From.ID, "failed to process")
		return
	}

	SendMessage(update.Message.From.ID, "done")
}
