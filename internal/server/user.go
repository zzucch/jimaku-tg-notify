package server

import (
	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

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
