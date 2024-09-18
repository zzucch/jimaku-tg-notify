package server

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

func (s *Server) Subscribe(chatID int64, titleID int64) error {
	client, err := s.clientManager.GetClient(chatID)
	if err != nil {
		log.Error(
			"failed to get client",
			"titleID",
			titleID,
			"err",
			err)

		return err
	}

	exists, err := storage.SubscriptionExists(chatID, titleID)
	if err != nil {
		log.Warn("failed to get subscription existence",
			"titleID",
			titleID,
			"err",
			err)
	}

	if exists {
		return errors.New("Already subscribed")
	}

	entry, err := client.GetEntryDetails(titleID)
	if err != nil {
		log.Warn("failed to get entry details",
			"titleID",
			titleID,
			"err",
			err)

		return err
	}

	latestSubtitleTimestamp, err := entry.GetLatestSubtitleTimestamp()
	if err != nil {
		log.Warn("failed to get latest subtitle timestamp",
			"titleID",
			titleID,
			"err",
			err)

		return err
	}

	if err := storage.Subscribe(
		chatID,
		titleID,
		latestSubtitleTimestamp,
		entry.JapaneseName,
	); err != nil {
		log.Warn(
			"failed to subscribe",
			"chatID",
			chatID,
			"titleID",
			titleID,
			"latestSubtitleTime",
			latestSubtitleTimestamp,
			"entry",
			entry,
			"err",
			err)

		return err
	}

	return nil
}

func (s *Server) Unsubscribe(chatID int64, titleID int64) error {
	if err := storage.Unsubscribe(chatID, titleID); err != nil {
		log.Warn(
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

func (s *Server) SetLatestTimestamp(
	chatID int64,
	titleID int64,
	latestTimestamp int64,
) {
	if err := storage.SetLatestSubtitleTimestamp(
		chatID,
		titleID,
		latestTimestamp,
	); err != nil {
		log.Error(
			"failed to set latest timestamp",
			"chatID",
			chatID,
			"titleID",
			titleID,
			"err",
			err)
	}
}

func (s *Server) SetJapaneseName(
	chatID int64,
	titleID int64,
	japaneseName string,
) {
	if err := storage.SetJapaneseName(
		chatID,
		titleID,
		japaneseName,
	); err != nil {
		log.Error(
			"failed to set japanese name",
			"chatID",
			chatID,
			"japanese name",
			japaneseName,
			"err",
			err)
	}
}
