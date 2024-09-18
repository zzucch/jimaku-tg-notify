package notification

import (
	"runtime"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
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
	store *storage.Storage,
) {
	go func() {
		notify(chatID, notificationCh, client, store)

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
				notify(chatID, notificationCh, client, store)
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

type SchedulerUpdate struct {
	ChatID   int64
	Interval time.Duration
}

func (m *Manager) WatchForSchedulerUpdates() {
	workerCount := runtime.NumCPU()

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for update := range m.updateCh {
				err := m.AddScheduler(update.ChatID, update.Interval)
				if err != nil {
					log.Error("failed to update scheduler", "update", update)

					m.notificationCh <- Notification{
						ChatID:  update.ChatID,
						Message: "Failed due to a critical error - contact the developers",
					}
				}
			}
		}()
	}

	wg.Wait()
}
