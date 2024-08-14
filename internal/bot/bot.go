package bot

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
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
	handleCommand(update, subscribeCommand)
}

func handleUnsubscribe(update tgbotapi.Update) {
	handleCommand(update, unsubscribeCommand)
}

func handleCommand(update tgbotapi.Update, command string) {
	chatID := update.Message.Chat.ID
	unvalidatedTitleID := update.Message.Text[len(command):]
	log.Debug("handling "+command, "chatID", chatID, "text", unvalidatedTitleID)

	unvalidatedTitleID = strings.TrimSpace(unvalidatedTitleID)
	_, err := strconv.Atoi(unvalidatedTitleID)
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

	fullURL := baseURL +
		command +
		"/" +
		url.PathEscape(strconv.FormatInt(chatID, 10)) +
		"/" +
		url.PathEscape(unvalidatedTitleID)

	resp, err := http.Get(fullURL)
	if err != nil {
		SendMessage(update.Message.From.ID, "failed to "+command)
		log.Error("failed to send a request", "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		SendMessage(update.Message.From.ID, "failed to "+command)
		log.Error("unexpected status code", "statusCode", resp.StatusCode)
		return
	}

	SendMessage(update.Message.From.ID, "done")
}
