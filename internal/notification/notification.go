package notification

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
	Updates []Update
}

type Update struct {
	TitleID         int64
	LatestTimestamp int64
}

func Notify(
	chatID int64,
	notificationCh chan Notification,
	client *client.Client,
) {
	var notificationMessageSB strings.Builder

	subscriptions, err := storage.GetAllSubscriptions(chatID)
	if err != nil {
		log.Error(
			"failed to get all subscriptions",
			"chatID",
			chatID,
			"err",
			err)

		notificationCh <- Notification{
			ChatID:  chatID,
			Message: "Failed due to a critical error - contact the developers",
		}
	}

	updates := make([]Update, 0, len(subscriptions))

	for _, subscription := range subscriptions {
		update, err := getUpdate(subscription.TitleID, client)
		if err != nil {
			log.Warn(
				"failed to get update",
				"titleID",
				subscription.TitleID,
				"err",
				err)
		}

		message := getUpdateMessage(subscription, update, err)
		notificationMessageSB.WriteString(message)

		if err == nil &&
			subscription.LatestSubtitleTime != update.LatestTimestamp {
			updates = append(updates, update)
		}
	}

	if notificationMessageSB.Len() == 0 {
		return
	}

	notificationCh <- Notification{
		ChatID:  chatID,
		Message: notificationMessageSB.String(),
		Updates: updates,
	}
}

func getUpdate(
	titleID int64,
	client *client.Client,
) (Update, error) {
	latestTimestamp, err := client.GetLatestSubtitleTime(titleID)
	if err != nil {
		return Update{}, err
	}

	return Update{
		TitleID:         titleID,
		LatestTimestamp: latestTimestamp,
	}, nil
}

func getUpdateMessage(
	subscription storage.Subscription,
	update Update,
	err error,
) string {
	var sb strings.Builder

	if err != nil {
		sb.WriteString("Failed to get latest subtitle date for jimaku.cc/entry/")
		sb.WriteString(strconv.FormatInt(subscription.TitleID, 10))
		sb.WriteString(":\n")
		sb.WriteString(err.Error())
		sb.WriteString("\n\n")
	} else if subscription.LatestSubtitleTime != update.LatestTimestamp {
		sb.WriteString("Update at jimaku.cc/entry/")
		sb.WriteString(strconv.FormatInt(subscription.TitleID, 10))
		sb.WriteString(" at ")
		sb.WriteString(util.TimestampToString(update.LatestTimestamp))
		sb.WriteString("\n")
	}

	return sb.String()
}
