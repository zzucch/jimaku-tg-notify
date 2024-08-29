package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleApiKeyChange(update tgbotapi.Update) {
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

	if err := b.server.SetApiKey(chatID, apiKey); err != nil {
		b.SendMessage(chatID, "Failed to process.\n"+err.Error())
		return
	}

	b.SendMessage(chatID, "done")
}