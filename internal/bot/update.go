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
	if update.LatestTimestamp != -1 && len(update.NewFileEntryNames) > 0 {
		if err := store.SetLatestSubtitleTimestamp(
			chatID,
			update.TitleID,
			update.LatestTimestamp,
		); err != nil {
			log.Error(
				"failed to set latest timestamp",
				"update",
				update,
				"err",
				err)
		}
	}

	if update.Name != "" {
		if err := store.SetName(
			chatID,
			update.TitleID,
			update.Name,
		); err != nil {
			log.Error(
				"failed to set name",
				"update",
				update,
				"err",
				err)
		}
	}
}
