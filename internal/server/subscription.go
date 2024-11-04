package server

import (
	"errors"
	"sort"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

func (s *Server) Subscribe(chatID int64, titleID int64) (string, error) {
	client, err := s.clientManager.GetClient(chatID)
	if err != nil {
		log.Error(
			"failed to get client",
			"titleID",
			titleID,
			"err",
			err,
		)

		return "", err
	}

	exists, err := s.store.SubscriptionExists(chatID, titleID)
	if err != nil {
		log.Warn("failed to get subscription existence",
			"titleID",
			titleID,
			"err",
			err,
		)
	}

	if exists {
		return "", errors.New("Already subscribed")
	}

	entry, err := client.GetEntryDetails(titleID)
	if err != nil {
		log.Warn("failed to get entry details",
			"titleID",
			titleID,
			"err",
			err,
		)

		return "", err
	}

	latestSubtitleTimestamp, err := entry.GetLastModified()
	if err != nil {
		log.Warn("failed to get latest subtitle timestamp",
			"titleID",
			titleID,
			"err",
			err,
		)

		return "", err
	}

	name := ""
	if entry.JapaneseName != "" {
		name = entry.JapaneseName
	} else if entry.Name != "" {
		name = entry.Name
	} else if entry.EnglishName != "" {
		name = entry.EnglishName
	}

	if name == "" {
		log.Error("failed to assign name value", "entry", entry)
	}

	if err := s.store.Subscribe(
		chatID,
		titleID,
		latestSubtitleTimestamp,
		name,
	); err != nil {
		log.Warn(
			"failed to subscribe",
			"titleID",
			titleID,
			"latestSubtitleTime",
			latestSubtitleTimestamp,
			"entry",
			entry,
			"err",
			err,
		)

		return "", err
	}

	return name, nil
}

func (s *Server) Unsubscribe(chatID int64, titleID int64) (string, error) {
	subscription, err := s.store.GetSubscription(chatID, titleID)
	if err != nil {
		log.Warn(
			"failed to find subscription",
			"titleID",
			titleID,
			"err",
			err,
		)
	}

	if err := s.store.Unsubscribe(chatID, titleID); err != nil {
		log.Warn(
			"failed to unsubscribe",
			"titleID",
			titleID,
			"err",
			err,
		)

		return "", err
	}

	return subscription.Name, nil
}

func (s *Server) ListSubscriptions(
	chatID int64,
) ([]storage.Subscription, error) {
	subscriptions, err := s.store.GetAllSubscriptions(chatID)
	if err != nil {
		log.Error(
			"failed to get all subscriptions",
			"err",
			err,
		)

		return nil, err
	}

	sort.Slice(subscriptions, func(i, j int) bool {
		return subscriptions[i].LastModified > subscriptions[j].LastModified
	})

	return subscriptions, nil
}

func (s *Server) SetLatestTimestamp(
	chatID int64,
	titleID int64,
	latestTimestamp int64,
) {
	if err := s.store.SetLatestSubtitleTimestamp(
		chatID,
		titleID,
		latestTimestamp,
	); err != nil {
		log.Error(
			"failed to set latest timestamp",
			"titleID",
			titleID,
			"err",
			err,
		)
	}
}

func (s *Server) SetName(
	chatID int64,
	titleID int64,
	name string,
) {
	if err := s.store.SetName(
		chatID,
		titleID,
		name,
	); err != nil {
		log.Error(
			"failed to set name",
			"name",
			name,
			"err",
			err,
		)
	}
}
