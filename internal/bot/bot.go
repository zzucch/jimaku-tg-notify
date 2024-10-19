package bot

import (
	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
	"github.com/zzucch/jimaku-tg-notify/internal/notification"
	"github.com/zzucch/jimaku-tg-notify/internal/server"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

type Bot struct {
	botAPI         *tgbotapi.BotAPI
	cache          loggedUsersCache
	server         *server.Server
	store          *storage.Storage
	notificationCh chan notification.Notification
}

const (
	listCommand        = "/list"
	subscribeCommand   = "/sub"
	unsubscribeCommand = "/unsub"
	apiKeyCommand      = "/apikey"
	intervalCommand    = "/interval"
	offsetCommand      = "/utc_offset"
)

func NewBot(
	config config.Config,
	server *server.Server,
	store *storage.Storage,
	notificationCh chan notification.Notification,
) (*Bot, error) {
	var err error

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		return &Bot{}, err
	}

	bot.Debug = config.BotDebugLevel

	return &Bot{
		botAPI:         bot,
		server:         server,
		store:          store,
		notificationCh: notificationCh,
	}, nil
}

func (b *Bot) Start() {
	log.Info("starting bot")

	go b.handleNotifications()

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
