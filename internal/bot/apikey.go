package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleAPIKeyChange(update tgbotapi.Update) {
	chatID := update.Message.From.ID

	apiKey := update.Message.Text[len(apiKeyCommand):]
	apiKey = strings.TrimSpace(apiKey)
	split := strings.Split(apiKey, " ")

	if apiKey == "" || len(split) > 1 {
		b.SendMessage(
			chatID,
			"Example usage:\n"+
				apiKeyCommand+
				" ZXhhbXBsZSBhcGkga2V5IGV4YW1wbGUg")

		return
	}

	if err := b.server.SetAPIKey(chatID, apiKey); err != nil {
		b.SendMessage(chatID, "Failed to process.\n"+err.Error())
		return
	}

	if !b.cache.exists(chatID) {
		b.cache.insert(chatID)
		b.handleHelp(update)
	} else {
		b.SendMessage(chatID, "Done")
	}
}
