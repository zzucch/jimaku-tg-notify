package rate

import (
	"sync"
	"time"
)

type Limiter struct {
	rwMutex        sync.RWMutex
	limit          int
	remaining      int
	resetTimestamp int64
	waitingAmount  int
}

func NewLimiter() *Limiter {
	return &Limiter{
		limit:          1,
		remaining:      1,
		resetTimestamp: time.Now().Add(100 * time.Second).Unix(),
	}
}

func (l *Limiter) SetLimit(limit int) {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	l.limit = limit
}

func (l *Limiter) SetRemaining(remaining int) {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	l.remaining = remaining
}

func (l *Limiter) SetResetTime(resetTime int64) {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	l.resetTimestamp = resetTime
}

func (l *Limiter) Wait() {
	l.rwMutex.Lock()
	l.waitingAmount++
	l.rwMutex.Unlock()

	defer func() {
		l.rwMutex.Lock()
		l.waitingAmount--
		l.rwMutex.Unlock()
	}()

	delay := l.getDelay()

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return
	}
}

func (l *Limiter) getDelay() time.Duration {
	l.rwMutex.RLock()
	defer l.rwMutex.RUnlock()

	if l.remaining > 0 {
		return time.Duration(0)
	}

	// HACK: assuming that if someone got into waiting with someone else then
	// both would make more then limit*2 requests,
	// probably would change in the future. it is not critical issue at all though
	return time.Duration(time.Until(
		time.Unix(l.resetTimestamp, 0)).Nanoseconds() * int64(l.waitingAmount))
}
