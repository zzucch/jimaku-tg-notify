package notification

import (
	"time"

	"github.com/zzucch/jimaku-tg-notify/internal/client"
)

type Command struct {
	ChatID      int64
	NewInterval time.Duration
}

type Scheduler struct {
	interval  time.Duration
	commandCh chan Command
	stopCh    chan struct{}
}

func NewScheduler(interval time.Duration) *Scheduler {
	return &Scheduler{
		interval:  interval,
		commandCh: make(chan Command),
		stopCh:    make(chan struct{}),
	}
}

func (s *Scheduler) Start(
	chatID int64,
	notificationCh chan Notification,
	client *client.Client,
) {
	go func() {
		Notify(chatID, notificationCh, client)

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-s.stopCh:
				return
			case cmd := <-s.commandCh:
				if cmd.ChatID == chatID {
					ticker.Stop()

					s.interval = cmd.NewInterval
					ticker = time.NewTicker(s.interval)
				}
			case <-ticker.C:
				Notify(chatID, notificationCh, client)
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}

func (s *Scheduler) UpdateInterval(
	chatID int64,
	newInterval time.Duration,
) {
	s.commandCh <- Command{ChatID: chatID, NewInterval: newInterval}
}
