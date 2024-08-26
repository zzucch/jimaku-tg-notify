package bot

import (
	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
)

var bot *tgbotapi.BotAPI

const (
	listCommand        = "/list"
	subscribeCommand   = "/sub"
	unsubscribeCommand = "/unsub"
)

func SendMessage(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(message); err != nil {
		log.Error("failed to send message", "err", err)
	}
}

func Initialize(config config.Config) error {
	log.Info("starting bot")
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
