package notify

import (
	"sync"
	"time"

	"github.com/zzucch/jimaku-tg-notify/internal/client"
)

type NotifyManager struct {
	schedulers     sync.Map
	notificationCh chan Notification
	client         *client.Client
}

func NewNotifyManager(
	notificationCh chan Notification,
	client *client.Client,
) *NotifyManager {
	return &NotifyManager{
		schedulers: sync.Map{},
		client:     client,
    notificationCh: notificationCh,
	}
}

func (nm *NotifyManager) AddScheduler(
	chatID int64,
	interval time.Duration,
) {
	if scheduler, exists := nm.schedulers.Load(chatID); exists {
		scheduler.(*NotifyScheduler).UpdateInterval(chatID, interval)
	} else {
		scheduler := NewNotifyScheduler(interval)

		nm.schedulers.Store(chatID, scheduler)
		scheduler.Start(chatID, nm.notificationCh, nm.client)
	}
}

func (nm *NotifyManager) RemoveScheduler(chatID int64) {
	if scheduler, exists := nm.schedulers.Load(chatID); exists {
		scheduler.(*NotifyScheduler).Stop()
		nm.schedulers.Delete(chatID)
	}
}

func (nm *NotifyManager) StopAll() {
	nm.schedulers.Range(func(key, value interface{}) bool {
		value.(*NotifyScheduler).Stop()
		nm.schedulers.Delete(key)
		return true
	})
}
