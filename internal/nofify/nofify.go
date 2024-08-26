package nofify

import (
	"log"
	"strconv"

	"github.com/zzucch/jimaku-tg-notify/internal/bot"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

func Notify() {
	chatIDs, err := storage.GetAllChatIDs()
	if err != nil {
		log.Fatal("failed to get all chat ids", "err", err)
	}

	for _, chatID := range chatIDs {
		subscriptions, err := storage.GetAllSubscriptions(chatID)
		if err != nil {
			log.Fatal("failed to get all subscriptions", "chatID", chatID, "err", err)
		}

		for _, subscription := range subscriptions {
			bot.SendMessage(
				chatID,
				strconv.FormatInt(subscription.TitleID, 10)+
					" "+
					strconv.FormatInt(subscription.LatestSubtitleTime, 10))
		}
	}
}
