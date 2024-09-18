package notification

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

type Notification struct {
	ChatID  int64
	Message string
	Updates []Update
}

type Update struct {
	TitleID           int64
	LatestTimestamp   int64
	JapaneseName      string
	NewFileEntryNames []string
}

func notify(
	chatID int64,
	notificationCh chan Notification,
	client *client.Client,
	store *storage.Storage,
) {
	var notificationMessageSB strings.Builder

	subscriptions, err := store.GetAllSubscriptions(chatID)
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
		update, err := getUpdate(subscription, client)
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

		if err == nil {
			if update.LatestTimestamp != 0 || update.JapaneseName != "" {
				updates = append(updates, update)
			}
		}
	}

	if notificationMessageSB.Len() == 0 {
		if len(updates) != 0 {
			notificationCh <- Notification{
				ChatID:  chatID,
				Message: "",
				Updates: updates,
			}
		}

		return
	}

	notificationCh <- Notification{
		ChatID:  chatID,
		Message: "Updates:\n" + notificationMessageSB.String(),
		Updates: updates,
	}
}
