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

func NotifyAll() {
	log.Info("notifying")
	chatIDs, err := storage.GetAllChatIDs()
	if err != nil {
		log.Fatal("failed to get all chat ids", "err", err)
	}

	for _, chatID := range chatIDs {
		var notificationMessageSB strings.Builder

		subscriptions, err := storage.GetAllSubscriptions(chatID)
		if err != nil {
			log.Fatal("failed to get all subscriptions", "chatID", chatID, "err", err)
		}

		for _, subscription := range subscriptions {
			notificationMessageSB.WriteString(getNotificationMessage(subscription))
		}

		if notificationMessageSB.Len() == 0 {
			continue
		}

		bot.SendMessage(chatID, notificationMessageSB.String())
	}
}

func getNotificationMessage(subscription storage.Subscription) string {
	latest, err := client.GetLatestSubtitleTimestamp(subscription.TitleID)
	if err != nil {
		log.Error(
			"failed to get latest subtitle timestamp",
			"subscription",
			subscription,
			"err",
			err)
	}

	if subscription.LatestSubtitleTime == latest {
		return ""
	}

	return "Update at jimaku.cc/entry/" +
		strconv.FormatInt(subscription.TitleID, 10) +
		" at time " +
		util.TimestampToString(subscription.LatestSubtitleTime) +
		"\n"
}

