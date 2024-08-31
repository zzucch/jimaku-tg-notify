package notification

import (
	"sync"
	"time"

	"github.com/zzucch/jimaku-tg-notify/internal/client"
)

type Manager struct {
	schedulers     sync.Map
	clientManager  *client.Manager
	updateCh       chan SchedulerUpdate
	notificationCh chan Notification
}

func NewManager(
	clientManager *client.Manager,
	updateCh chan SchedulerUpdate,
	notificationCh chan Notification,
) *Manager {
	return &Manager{
		clientManager:  clientManager,
		updateCh:       updateCh,
		notificationCh: notificationCh,
	}
}

func (m *Manager) AddScheduler(
	chatID int64,
	interval time.Duration,
) error {
	if scheduler, exists := m.schedulers.Load(chatID); exists {
		scheduler.(*Scheduler).UpdateInterval(chatID, interval)

		return nil
	}

	scheduler := NewScheduler(interval)
	m.schedulers.Store(chatID, scheduler)

	client, err := m.clientManager.GetClient(chatID)
	if err != nil {
		return err
	}

	scheduler.Start(chatID, m.notificationCh, client)

	return nil
}

func (m *Manager) RemoveScheduler(chatID int64) {
	if scheduler, exists := m.schedulers.Load(chatID); exists {
		scheduler.(*Scheduler).Stop()
		m.schedulers.Delete(chatID)
	}
}

func (m *Manager) StopAll() {
	m.schedulers.Range(func(key, value interface{}) bool {
		value.(*Scheduler).Stop()
		m.schedulers.Delete(key)

		return true
	})
}
