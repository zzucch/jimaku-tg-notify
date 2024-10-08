package client

import (
	"sync"

	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

type Manager struct {
	clients sync.Map
	store   *storage.Storage
}

func NewManager(store *storage.Storage) *Manager {
	return &Manager{
		store: store,
	}
}

func (m *Manager) GetClient(chatID int64) (*Client, error) {
	unvalidated, ok := m.clients.Load(chatID)
	if !ok {
		apiKey, err := m.store.GetAPIKey(chatID)
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
	apiKey, err := m.store.GetAPIKey(chatID)
	if err != nil {
		return err
	}

	unvalidated, ok := m.clients.Load(chatID)
	if ok {
		client := unvalidated.(*Client)
		client.UpdateAPIKey(apiKey)
	} else {
		c := NewClient(apiKey)
		m.clients.Store(chatID, c)
	}

	return nil
}
