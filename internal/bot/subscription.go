package bot

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleSubscribe(update tgbotapi.Update) {
	b.handleSubscription(
		update,
		subscribeCommand,
		b.server.Subscribe,
		"Subscribed to ",
	)
}

func (b *Bot) handleUnsubscribe(update tgbotapi.Update) {
	b.handleSubscription(
		update,
		unsubscribeCommand,
		b.server.Unsubscribe,
		"Unsubscribed from ",
	)
}

func (b *Bot) handleSubscription(
	update tgbotapi.Update,
	command string,
	action func(chatID int64, titleID int64) (string, error),
	doneMessage string,
) {
	chatID := update.Message.Chat.ID

	unvalidatedTitleID := update.Message.Text[len(command):]
	unvalidatedTitleID = strings.TrimSpace(unvalidatedTitleID)

	titleID, err := strconv.ParseInt(unvalidatedTitleID, 10, 64)
	if unvalidatedTitleID == "" || err != nil || titleID < 0 {
		if err := b.SendMessage(
			chatID,
			"Example usage:\n"+command+" 123",
		); err != nil {
			log.Error("failed to send message", "err", err)
		}

		return
	}

	name, err := action(chatID, titleID)
	if err != nil {
		if err := b.SendMessage(
			chatID,
			"Failed to process.\n"+err.Error(),
		); err != nil {
			log.Error("failed to send message", "err", err)
		}

		return
	}

	if err := b.SendMessage(
		update.Message.From.ID,
		doneMessage+name,
	); err != nil {
		log.Error("failed to send message", "err", err)
	}
}
