package notify

import (
	"runtime"
	"sync"
	"time"
)

type Update struct {
	ChatID   int64
	Interval time.Duration
}

func (nm *NotifyManager) WatchForUpdates() {
	workerCount := runtime.NumCPU()
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for update := range nm.updateCh {
				nm.AddScheduler(update.ChatID, update.Interval)
			}
		}()
	}

	wg.Wait()
}
