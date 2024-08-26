package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/server"
	"github.com/zzucch/jimaku-tg-notify/internal/util"
)

func handleSubscriptionList(update tgbotapi.Update) {
	chatID := update.Message.From.ID

	var messageSB strings.Builder
	subscriptions, err := server.ListSubscriptions(chatID)
	if err != nil {
		SendMessage(chatID, "failed to process")
	}

	if len(subscriptions) == 0 {
		messageSB.WriteString("You don't have any subscriptions yet! ")
		messageSB.WriteString("To subscribe, use ")
		messageSB.WriteString(subscribeCommand)
		messageSB.WriteString(" [title_id]")
	} else {
		messageSB.WriteString("Subscriptions list (entry - last update):")
	}

	for _, subscription := range subscriptions {
		messageSB.WriteString("\njimaku.cc/entry/")
		messageSB.WriteString(strconv.FormatInt(subscription.TitleID, 10))
		messageSB.WriteString(" - ")
		messageSB.WriteString(
			util.TimestampToString(subscription.LatestSubtitleTime))
	}

	SendMessage(chatID, messageSB.String())
}
