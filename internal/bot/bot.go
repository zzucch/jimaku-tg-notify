package bot

import (
	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
	"github.com/zzucch/jimaku-tg-notify/internal/server"
)

type Bot struct {
	botAPI *tgbotapi.BotAPI
	server server.Server
}

const (
	listCommand        = "/list"
	subscribeCommand   = "/sub"
	unsubscribeCommand = "/unsub"
)

func (b *Bot) SendMessage(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)
	if _, err := b.botAPI.Send(message); err != nil {
		log.Error("failed to send message", "err", err)
	}
}

func Initialize(config config.Config) (Bot, error) {
	log.Info("starting bot")
	var err error
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		return Bot{}, err
	}

	if config.BotDebugLevel {
		bot.Debug = true
	}

	return Bot{bot, server.Server{*client.NewClient(config.APIKey)}}, nil
}

func (b *Bot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.botAPI.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		b.handleMessage(update)
	}
}
