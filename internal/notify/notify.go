package notify

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/bot"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
	"github.com/zzucch/jimaku-tg-notify/internal/util"
)

func Notify(
	chatID int64,
	bot *bot.Bot,
	client *client.Client,
) {
	var notificationMessageSB strings.Builder

	subscriptions, err := storage.GetAllSubscriptions(chatID)
	if err != nil {
		bot.SendMessage(
			chatID,
			"failed due to critical error - contact developers")

		log.Fatal(
			"failed to get all subscriptions",
			"chatID",
			chatID,
			"err",
			err)
	}

	for _, subscription := range subscriptions {
		notificationMessageSB.WriteString(
			getNotificationMessage(subscription, client))
	}

	if notificationMessageSB.Len() == 0 {
		return
	}

	bot.SendMessage(chatID, notificationMessageSB.String())
}

func getNotificationMessage(
	subscription storage.Subscription,
	client *client.Client,
) string {
	latestSubtitleTime, err := client.GetLatestSubtitle(subscription.TitleID)
	if err != nil {
		log.Error("failed to get latest subtitle date",
			"titleID",
			subscription.TitleID,
			"err",
			err)

		return "failed to get latest subtitle date"
	}

	if subscription.LatestSubtitleTime == latestSubtitleTime {
		return ""
	}

	return "Update at jimaku.cc/entry/" +
		strconv.FormatInt(subscription.TitleID, 10) +
		" at time " +
		util.TimestampToString(subscription.LatestSubtitleTime) +
		"\n"
}