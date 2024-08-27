package notify

import (
	"time"

	"github.com/zzucch/jimaku-tg-notify/internal/bot"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
)

type Command struct {
	ChatID      int64
	NewInterval time.Duration
}

type NotifyScheduler struct {
	interval  time.Duration
	commandCh chan Command
	stopCh    chan struct{}
}

func NewNotifyScheduler(interval time.Duration) *NotifyScheduler {
	return &NotifyScheduler{
		interval:  interval,
		commandCh: make(chan Command),
		stopCh:    make(chan struct{}),
	}
}

func (ns *NotifyScheduler) Start(
	chatID int64,
	bot *bot.Bot,
	client *client.Client,
) {
	go func() {
		ticker := time.NewTicker(ns.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				Notify(chatID, bot, client)
			case cmd := <-ns.commandCh:
				if cmd.ChatID == chatID {
					ticker.Stop()
					ticker = time.NewTicker(cmd.NewInterval)
				}
			case <-ns.stopCh:
				return
			}
		}
	}()
}

func (ns *NotifyScheduler) Stop() {
	close(ns.stopCh)
}

func (ns *NotifyScheduler) UpdateInterval(
	chatID int64,
	newInterval time.Duration,
) {
	ns.commandCh <- Command{ChatID: chatID, NewInterval: newInterval}
}