package server

import (
	"sync"

	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/notify"
)

type Server struct {
	users         sync.Map
	clientManager *client.ClientManager
	updateCh      chan notify.Update
}

func NewServer(
	chatIDs []int64,
	clientManager *client.ClientManager,
	updateCh chan notify.Update,
) *Server {
	server := &Server{
		users:         sync.Map{},
		clientManager: clientManager,
		updateCh:      updateCh,
	}

	for _, chatID := range chatIDs {
		server.users.LoadOrStore(chatID, struct{}{})
	}

	return server
}
