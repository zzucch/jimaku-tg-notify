package client

import (
	"sync"

	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

type Manager struct {
	clients sync.Map
}

func (m *Manager) GetClient(chatID int64) (*Client, error) {
	unvalidated, ok := m.clients.Load(chatID)
	if !ok {
		apiKey, err := storage.GetAPIKey(chatID)
		if err != nil {
			return nil, err
		}

		c := NewClient(apiKey)
		m.clients.Store(chatID, c)

		return c, nil
	}

	return unvalidated.(*Client), nil
}

func (m *Manager) UpdateAPIKey(chatID int64) error {
	apiKey, err := storage.GetAPIKey(chatID)
	if err != nil {
		return err
	}

	c := NewClient(apiKey)
	m.clients.Store(chatID, c)

	return nil
}
