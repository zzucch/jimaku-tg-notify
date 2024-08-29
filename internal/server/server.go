package server

import (
	"sync"

	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

type Server struct {
	clients sync.Map
	users   sync.Map
}

func NewServer(chatIDs []int64) *Server {
	server := &Server{
		clients: sync.Map{},
		users:   sync.Map{},
	}

	for _, chatID := range chatIDs {
		server.users.LoadOrStore(chatID, struct{}{})
	}

	return server
}

func (s *Server) getClient(chatID int64) (*client.Client, error) {
	v, ok := s.clients.Load(chatID)
	if !ok {
		apiKey, err := storage.GetApiKey(chatID)
		if err != nil {
			return nil, err
		}

		c := client.NewClient(apiKey)
		s.clients.Store(chatID, c)

		return c, err
	}

	return v.(*client.Client), nil
}
