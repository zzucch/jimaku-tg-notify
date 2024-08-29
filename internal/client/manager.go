package client

import (
	"sync"

	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

type ClientManager struct {
	clients sync.Map
}

func (cm *ClientManager) GetClient(chatID int64) (*Client, error) {
	v, ok := cm.clients.Load(chatID)
	if !ok {
		apiKey, err := storage.GetApiKey(chatID)
		if err != nil {
			return nil, err
		}

		c := NewClient(apiKey)
		cm.clients.Store(chatID, c)

		return c, err
	}

	return v.(*Client), nil
}
