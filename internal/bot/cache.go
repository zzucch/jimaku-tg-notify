package bot

import (
	"sync"
	"time"
)

type loggedUsersCache struct {
	sync.Map
}

func (c *loggedUsersCache) insert(chatID int64) {
	c.Store(chatID, struct{}{})

	expiration := time.After(6 * time.Hour)
	go func() {
		<-expiration
		c.Delete(chatID)
	}()
}

func (c *loggedUsersCache) exists(chatID int64) bool {
	_, ok := c.Load(chatID)

	return ok
}
