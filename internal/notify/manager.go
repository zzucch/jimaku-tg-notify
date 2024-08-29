package notify

import (
	"sync"
	"time"

	"github.com/zzucch/jimaku-tg-notify/internal/client"
)

type NotifyManager struct {
	schedulers     sync.Map
	notificationCh chan Notification
	clientManager  *client.ClientManager
}

func NewNotifyManager(
	notificationCh chan Notification,
	clientManager *client.ClientManager,
) *NotifyManager {
	return &NotifyManager{
		schedulers:     sync.Map{},
		clientManager:  clientManager,
		notificationCh: notificationCh,
	}
}

func (nm *NotifyManager) AddScheduler(
	chatID int64,
	interval time.Duration,
) error {
	if scheduler, exists := nm.schedulers.Load(chatID); exists {
		scheduler.(*NotifyScheduler).UpdateInterval(chatID, interval)

		return nil
	}

	scheduler := NewNotifyScheduler(interval)
	nm.schedulers.Store(chatID, scheduler)

	client, err := nm.clientManager.GetClient(chatID)
	if err != nil {
		return err
	}

	scheduler.Start(chatID, nm.notificationCh, client)

	return nil
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
