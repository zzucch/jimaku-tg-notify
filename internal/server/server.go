package server

import (
	"sync"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

type Server struct {
	client *client.Client
	users  sync.Map
}

func NewServer(chatIDs []int64, client *client.Client) *Server {
	server := &Server{
		client: client,
		users:  sync.Map{},
	}

	for _, chatID := range chatIDs {
		server.users.LoadOrStore(chatID, struct{}{})
	}

	return server
}

func (s *Server) AddUser(chatID int64) error {
	_, ok := s.users.LoadOrStore(chatID, struct{}{})
	if ok {
		return nil
	}

	err := storage.AddUser(chatID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Subscribe(chatID int64, titleID int64) error {
	latestSubtitleTime, err := s.client.GetLatestSubtitle(titleID)
	if err != nil {
		log.Error("failed to get latest subtitle date",
			"titleID",
			titleID,
			"err",
			err)

		return err
	}

	if err := storage.Subscribe(
		chatID,
		titleID,
		latestSubtitleTime); err != nil {
		log.Error(
			"failed to subscribe",
			"chatID",
			chatID,
			"titleID",
			titleID,
			"latestSubtitleTime",
			latestSubtitleTime,
			"err",
			err)

		return err
	}

	return nil
}

func (s *Server) Unsubscribe(chatID int64, titleID int64) error {
	if err := storage.Unsubscribe(chatID, titleID); err != nil {
		log.Error(
			"failed to unsubscribe",
			"chatID",
			chatID,
			"titleID",
			titleID,
			"err",
			err)

		return err
	}

	return nil
}

func (s *Server) ListSubscriptions(
	chatID int64,
) ([]storage.Subscription, error) {
	subscriptions, err := storage.GetAllSubscriptions(chatID)
	if err != nil {
		log.Error(
			"failed to get all subscriptions",
			"chatID",
			chatID,
			"err",
			err)

		return nil, err
	}

	return subscriptions, nil
}

func (s *Server) SetInterval(
	chatID int64,
	interval int,
) error {
	err := storage.SetNotificationInterval(chatID, interval)
	if err != nil {
		log.Error(
			"failed to set interval",
			"chatID",
			chatID,
			"err",
			err)

		return err
	}

	return nil
}
