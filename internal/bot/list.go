package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/timeutil"
)

func (b *Bot) handleSubscriptionList(update tgbotapi.Update) {
	chatID := update.Message.From.ID

	subscriptions, err := b.server.ListSubscriptions(chatID)
	if err != nil {
		_ = b.SendMessage(chatID, "Failed to process.\n"+err.Error())
	}

	var messageSB strings.Builder

	if len(subscriptions) == 0 {
		messageSB.WriteString("You don't have any subscriptions yet!\n")
		messageSB.WriteString("To subscribe, use ")
		messageSB.WriteString(subscribeCommand)
		messageSB.WriteString(" [title_id]")
	} else {
		messageSB.WriteString("Subscriptions list (entry, last update):")

		for _, subscription := range subscriptions {
			messageSB.WriteString("\n\n")
			messageSB.WriteString(subscription.JapaneseName)
			messageSB.WriteString("\njimaku.cc/entry/")
			messageSB.WriteString(strconv.FormatInt(subscription.TitleID, 10))
			messageSB.WriteString(" - ")
			messageSB.WriteString(
				timeutil.TimestampToString(subscription.LastModified))
		}

		messageSB.WriteString("\n\n")

		if lastUpdateCheckTimestamp, err := b.server.GetLastUpdateCheck(
			chatID,
		); err != nil {
			messageSB.WriteString("Failed to get last update check time")
		} else {
			messageSB.WriteString("Last updates check time:\n")
			messageSB.WriteString(
				timeutil.TimestampToString(lastUpdateCheckTimestamp))
		}
	}

	_ = b.SendMessage(chatID, messageSB.String())
}
