package bot

import (
	"strings"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleAPIKeyChange(update tgbotapi.Update) {
	chatID := update.Message.From.ID

	apiKey := update.Message.Text[len(apiKeyCommand):]
	apiKey = strings.TrimSpace(apiKey)
	split := strings.Split(apiKey, " ")

	if apiKey == "" || len(split) > 1 {
		if err := b.SendMessage(
			chatID,
			"Example usage:\n"+
				apiKeyCommand+
				" ZXhhbXBsZSBhcGkga2V5IGV4YW1wbGUg",
		); err != nil {
			log.Error("failed to send message", "err", err)
		}

		return
	}

	if err := b.server.SetAPIKey(chatID, apiKey); err != nil {
		if err := b.SendMessage(
			chatID,
			"Failed to process.\n"+err.Error(),
		); err != nil {
			log.Error("failed to send message", "err", err)
		}
		return
	}

	if !b.cache.exists(chatID) {
		b.cache.insert(chatID)
		b.handleHelp(update)
	} else {
		if err := b.SendMessage(chatID, "API key is set"); err != nil {
			log.Error("failed to send message", "err", err)
		}
	}
}
