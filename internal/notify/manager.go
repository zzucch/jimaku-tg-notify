package notify

import (
	"sync"
	"time"

	"github.com/zzucch/jimaku-tg-notify/internal/bot"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
)

type NotifyManager struct {
	schedulers sync.Map
	bot        *bot.Bot
	client     *client.Client
}

func NewNotifyManager(bot *bot.Bot, client *client.Client) *NotifyManager {
	return &NotifyManager{
		schedulers: sync.Map{},
		bot:        bot,
		client:     client,
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
		scheduler.Start(chatID, nm.bot, nm.client)
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
