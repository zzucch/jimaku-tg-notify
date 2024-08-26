package server

import (
	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

func Subscribe(chatID int64, titleID int64) error {
	latestSubtitleTime, err := client.GetLatestSubtitleTimestamp(titleID)
	if err != nil {
		log.Error("failed to get subtitle dates", "titleID", titleID, "err", err)
		return err
	}

	if err := storage.Subscribe(chatID, titleID, latestSubtitleTime); err != nil {
		log.Error(
			"failed to subscribe",
			"chatID",
			chatID,
			"titleID",
			titleID,
			"latestSubtitleTime",
			latestSubtitleTime,
			"err",
			err)
		return err
	}

	return nil
}

func Unsubscribe(chatID int64, titleID int64) error {
  if err := storage.Unsubscribe(chatID, titleID); err != nil {
		log.Error(
			"failed to unsubscribe",
			"chatID",
			chatID,
			"titleID",
			titleID,
			"err",
			err)
		return err
	}

	return nil
}
