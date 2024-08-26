package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleHelp(update tgbotapi.Update) {
	helpMessage :=
		`
Available commands:
/list 
List all subscriptions

/sub [title_id]
Subscribe to all updates on given title

/unsub [title_id]
Unsubscribe from given title
`

	SendMessage(update.Message.From.ID, helpMessage)
}
