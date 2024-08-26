package bot

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
	"github.com/zzucch/jimaku-tg-notify/internal/server"
)

var bot *tgbotapi.BotAPI

const (
	subscribeCommand   = "/sub"
	unsubscribeCommand = "/unsub"
	baseURL            = "http://localhost:3001"
)

func SendMessage(chatID int64, text string) {
	log.Debug("sending message", "chatID", chatID, "messageText", text)

	message := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(message); err != nil {
		log.Error("failed to send message", "err", err)
	}
}

func Initialize(config config.Config) error {
	log.Debug("starting bot")
	var err error
	bot, err = tgbotapi.NewBotAPI(config.BotToken)
  if err != nil {
    return err
  }

	if config.BotDebugLevel {
		bot.Debug = true
	}

  return nil
}

func Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handleMessage(update)
	}
}

func handleMessage(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	log.Debug("handling message", "chatID", chatID)

	messageText := update.Message.Text
	command := strings.Split(messageText, " ")[0]

	switch command {
	case subscribeCommand:
		handleSubscribe(update)
	case unsubscribeCommand:
		handleUnsubscribe(update)
	default:
		log.Debug(
			"cannot handle message",
			"chatID",
			chatID,
			"update",
			update)
	}
}

func handleSubscribe(update tgbotapi.Update) {
	handleSubscription(update, subscribeCommand, server.Subscribe)
}

func handleUnsubscribe(update tgbotapi.Update) {
	handleSubscription(update, unsubscribeCommand, server.Unsubscribe)
}

func handleSubscription(
	update tgbotapi.Update,
	command string,
	action func(chatID int64, titleID int64) error,
) {
	chatID := update.Message.Chat.ID
	unvalidatedTitleID := update.Message.Text[len(command):]
	log.Debug("handling "+command, "chatID", chatID, "text", unvalidatedTitleID)

	unvalidatedTitleID = strings.TrimSpace(unvalidatedTitleID)
	titleID, err := strconv.ParseInt(unvalidatedTitleID, 10, 64)
	if unvalidatedTitleID == "" || err != nil {
		log.Debug(
			"failed to handle - invalid titleID",
			"titleID",
			unvalidatedTitleID,
			"err",
			err)

		SendMessage(chatID, "invalid command")
		return
	}

	if err := action(chatID, titleID); err != nil {
		SendMessage(update.Message.From.ID, "failed to process request")
		return
	}

	SendMessage(update.Message.From.ID, "done")
}
