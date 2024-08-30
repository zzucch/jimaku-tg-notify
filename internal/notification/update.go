package notification

import (
	"runtime"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type Update struct {
	ChatID   int64
	Interval time.Duration
}

func (m *Manager) WatchForUpdates() {
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
