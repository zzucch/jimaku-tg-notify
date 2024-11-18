package bot

import (
	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleLogin(update tgbotapi.Update) {
	message := `
Please set your API key:
/apikey [key]

If you have not generated an API key yet, you can do so on your account page: jimaku.cc/login
`

	if err := b.SendMessage(update.Message.From.ID, message); err != nil {
		log.Error("failed to send message", "err", err)
	}
}
