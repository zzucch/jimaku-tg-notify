package bot

import (
	"runtime"
	"sync"

	"github.com/charmbracelet/log"
)

func (b *Bot) handleNotifications() {
	workerCount := runtime.NumCPU()

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for notification := range b.notificationCh {
				if notification.Message != "" {
					if err := b.SendMessage(
						notification.ChatID,
						notification.Message,
					); err != nil {
						log.Error(
							"failed to send update message",
							"ignoring updates",
							len(notification.Updates))
					}
				}

				log.Debug("lk", "l", notification.Updates)
				if len(notification.Updates) > 0 {
					for _, update := range notification.Updates {
						updateStorage(b.store, notification.ChatID, update)
					}
				}
			}
		}()
	}

	wg.Wait()
}
