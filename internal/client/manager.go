package client

import (
	"log"
	"sync"

	"github.com/zzucch/jimaku-tg-notify/internal/storage"
	"github.com/zzucch/jimaku-tg-notify/pkg/client"
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

func (m *Manager) GetClient(chatID int64) (*client.Client, error) {
	unvalidated, ok := m.clients.Load(chatID)
	if !ok {
		apiKey, err := m.store.GetAPIKey(chatID)
		if err != nil {
			return nil, err
		}

		c := client.NewClient(apiKey)
		m.clients.Store(chatID, c)

		return c, nil
	}

	return unvalidated.(*client.Client), nil
}

func (m *Manager) UpdateAPIKey(chatID int64) error {
	apiKey, err := m.store.GetAPIKey(chatID)
	if err != nil {
		return err
	}

	unvalidated, ok := m.clients.Load(chatID)
	if ok {
		client, ok := unvalidated.(*client.Client)
		if !ok {
			log.Fatal(
				"invalid type",
				"expected",
				"*client.Client",
				"actual",
				unvalidated,
			)
		}

		client.UpdateAPIKey(apiKey)
	} else {
		c := client.NewClient(apiKey)
		m.clients.Store(chatID, c)
	}

	return nil
}
