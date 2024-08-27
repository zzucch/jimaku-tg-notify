package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleHelp(update tgbotapi.Update) {
	helpMessage :=
		`
Available commands:
/list 
List all subscriptions

/sub [title_id]
Subscribe to all updates on the given title

/unsub [title_id]
Unsubscribe from the given title

/interval [hours]
Set current notification interval to the given amount of hours
Default value is 6 hours
`

	b.SendMessage(update.Message.From.ID, helpMessage)
}
