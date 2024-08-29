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
Set current notification interval to the given amount of hours.
By default it is 6 hours

/apikey [key]
Set personal api key


How to get Title ID:
For example, the jimaku entry for 「逃げるは恥だが役に立つ」 is https://jimaku.cc/entry/3331, so the Title ID would be 3331
`

	b.SendMessage(update.Message.From.ID, helpMessage)
}
