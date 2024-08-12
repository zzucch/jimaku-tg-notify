package bot

import (
	"strings"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
)

var bot *tgbotapi.BotAPI

const (
	subscribeCommand   = "/sub"
	unsubscribeCommand = "/unsub"
)

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
	chatID := update.Message.Chat.ID
	text := update.Message.Text[len(subscribeCommand):]
	log.Debug("handling "+subscribeCommand, "chatID", chatID, "text", text)
}

func handleUnsubscribe(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text[len(unsubscribeCommand):]
	log.Debug("handling "+unsubscribeCommand, "chatID", chatID, "text", text)
}

func SendMessage(chatID int64, text string) {
	log.Debug("sending message", "chatID", chatID, "messageText", text)

	message := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(message); err != nil {
		log.Fatal(err)
	}
}

func Start(config config.Config) {
	log.Debug("starting bot")
	var err error
	bot, err = tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal(
			"failed creating BotAPI instance",
			"err",
			err,
			"token",
			config.BotConfig.BotToken)
	}

	if config.BotDebugLevel {
		bot.Debug = true
	}

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
