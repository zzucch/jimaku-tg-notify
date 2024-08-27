package server

import (
	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

type Server struct {
	client.Client
}

func (s *Server) Subscribe(chatID int64, titleID int64) error {
	latestSubtitleTime, err := s.Client.GetLatestSubtitle(titleID)
	if err != nil {
		log.Error("failed to get latest subtitle date",
			"titleID",
			titleID,
			"err",
			err)

		return err
	}

	err = storage.AddUser(chatID)
	if err != nil {
		log.Debug("failed to add user", "err", err)
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
