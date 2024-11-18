package bot

import (
	"strings"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(update tgbotapi.Update) {
	messageText := update.Message.Text
	command := strings.Split(messageText, " ")[0]

	chatID := update.Message.From.ID

	if serverErr := b.server.AddUser(update.Message.From.ID); serverErr != nil {
		if sendErr := b.SendMessage(
			chatID,
			"Failed due to a critical error - contact the developers",
		); sendErr != nil {
			log.Error("failed to send message", "err", sendErr)
		}

		log.Error(
			"failed to add user",
			"update",
			update,
			"messageText",
			messageText,
		)

		return
	}

	if !b.cache.exists(chatID) {
		if exists, validationErr := b.server.ValidateAPIKey(
			chatID,
		); validationErr != nil {
			if sendErr := b.SendMessage(
				chatID,
				"Failed due to a critical error - contact the developers",
			); sendErr != nil {
				log.Error("failed to send message", "err", sendErr)
			}

			return
		} else if exists {
			b.cache.insert(chatID)
		} else {
			b.handleUnauthenticatedCommand(command, update)

			return
		}
	}

	b.handleCommand(command, update)
}

func (b *Bot) handleUnauthenticatedCommand(
	command string,
	update tgbotapi.Update,
) {
	switch command {
	case apiKeyCommand:
		b.handleAPIKeyChange(update)
	default:
		b.handleLogin(update)
	}
}

func (b *Bot) handleCommand(command string, update tgbotapi.Update) {
	switch command {
	case subscribeCommand:
		b.handleSubscribe(update)
	case unsubscribeCommand:
		b.handleUnsubscribe(update)
	case listCommand:
		b.handleSubscriptionList(update)
	case apiKeyCommand:
		b.handleAPIKeyChange(update)
	case intervalCommand:
		b.handleNotificationIntervalChange(update)
	case offsetCommand:
		b.handleSettingUTCOffset(update)
	default:
		b.handleHelp(update)
	}
}
