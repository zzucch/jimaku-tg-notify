package notify

import "time"

type Update struct {
	ChatID   int64
	Interval time.Duration
}

func (nm *NotifyManager) WatchForUpdates() {
	for update := range nm.updateCh {
		nm.AddScheduler(update.ChatID, update.Interval)
	}
}
