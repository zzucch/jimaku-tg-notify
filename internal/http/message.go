package http

import (
	"encoding/json"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/zzucch/jimaku-tg-notify/internal/bot"
)

type SendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("handling sendMessage", "request", r)

	if r.Method != http.MethodPost {
		errorText := "Invalid request method"
		log.Debug("failed to handle sendMessage", "error", errorText)
		http.Error(w, errorText, http.StatusMethodNotAllowed)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorText := "Bad request"
		log.Debug("failed to handle sendMessage", "error", errorText)
		http.Error(w, errorText, http.StatusMethodNotAllowed)
		return
	}

	if req.ChatID == 0 || req.Text == "" {
		errorText := "ChatID and Text are required"
		log.Debug("failed to handle sendMessage", "error", errorText)
		http.Error(w, errorText, http.StatusMethodNotAllowed)
		return
	}

	bot.SendMessage(req.ChatID, req.Text)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message sent"))
	log.Debug("message sent")
}

func Start() {
	http.HandleFunc("/sendMessage", sendMessageHandler)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("failed to start server", "error", err)
		}
	}()

}
