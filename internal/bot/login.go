package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleLogin(update tgbotapi.Update) {
	message := `
Please set your API key:
/apikey [key]

If you have not generated an API key yet, you can do so on your account page: jimaku.cc/login
`

	b.SendMessage(update.Message.From.ID, message)
}
