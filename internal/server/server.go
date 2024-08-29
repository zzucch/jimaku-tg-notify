package server

import (
	"sync"

	"github.com/zzucch/jimaku-tg-notify/internal/client"
)

type Server struct {
	users         sync.Map
	clientManager *client.ClientManager
}

func NewServer(
	chatIDs []int64,
	clientManager *client.ClientManager,
) *Server {
	server := &Server{
		users:         sync.Map{},
		clientManager: clientManager,
	}

	for _, chatID := range chatIDs {
		server.users.LoadOrStore(chatID, struct{}{})
	}

	return server
}
