package server

import (
	"sync"

	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/notification"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

type Server struct {
	users         sync.Map
	store         *storage.Storage
	clientManager *client.Manager
	updateCh      chan notification.SchedulerUpdate
}

func NewServer(
	chatIDs []int64,
	store *storage.Storage,
	clientManager *client.Manager,
	updateCh chan notification.SchedulerUpdate,
) *Server {
	server := &Server{
		users:         sync.Map{},
		store:         store,
		clientManager: clientManager,
		updateCh:      updateCh,
	}

	for _, chatID := range chatIDs {
		server.users.LoadOrStore(chatID, struct{}{})
	}

	return server
}
