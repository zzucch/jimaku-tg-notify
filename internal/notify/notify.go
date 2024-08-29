package notify

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
	"github.com/zzucch/jimaku-tg-notify/internal/util"
)

type Notification struct {
	ChatID  int64
	Message string
}

func Notify(
	chatID int64,
	notificationCh chan Notification,
	client *client.Client,
) {
	var notificationMessageSB strings.Builder

	subscriptions, err := storage.GetAllSubscriptions(chatID)
	if err != nil {
		notificationCh <- Notification{
			ChatID:  chatID,
			Message: "Failed due to a critical error - contact the developers",
		}

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

	notificationCh <- Notification{
		ChatID:  chatID,
		Message: notificationMessageSB.String(),
	}
}

func getNotificationMessage(
	subscription storage.Subscription,
	client *client.Client,
) string {
	latestSubtitleTime, err := client.GetLatestSubtitleTime(subscription.TitleID)
	if err != nil {
		log.Error("failed to get latest subtitle date",
			"titleID",
			subscription.TitleID,
			"err",
			err)

		var sb strings.Builder
		sb.WriteString("Failed to get latest subtitle date for jimaku.cc/entry/")
		sb.WriteString(strconv.FormatInt(subscription.TitleID, 10))
		sb.WriteString(":\n")
		sb.WriteString(err.Error())
		sb.WriteString("\n\n")

		return sb.String()
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
