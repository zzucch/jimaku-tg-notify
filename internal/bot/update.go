package bot

import (
	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/notification"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

func updateStorage(
	store *storage.Storage,
	chatID int64,
	update notification.Update,
) {
	if update.LatestTimestamp != -1 {
		if err := store.SetLatestSubtitleTimestamp(
			chatID,
			update.TitleID,
			update.LatestTimestamp,
		); err != nil {
			log.Error(
				"failed to set latest timestamp",
				"chatID",
				chatID,
				"update",
				update,
				"err",
				err)
		}
	}

	if update.JapaneseName != "" {
		if err := store.SetJapaneseName(
			chatID,
			update.TitleID,
			update.JapaneseName,
		); err != nil {
			log.Error(
				"failed to set japanese name",
				"chatID",
				chatID,
				"update",
				update,
				"err",
				err)
		}
	}
}