package bot

import (
	"runtime"
	"sync"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zzucch/jimaku-tg-notify/internal/config"
	"github.com/zzucch/jimaku-tg-notify/internal/notify"
	"github.com/zzucch/jimaku-tg-notify/internal/server"
)

type Bot struct {
	botAPI         *tgbotapi.BotAPI
	cache          loggedUsersCache
	server         *server.Server
	notificationCh chan notify.Notification
}

const (
	listCommand        = "/list"
	subscribeCommand   = "/sub"
	unsubscribeCommand = "/unsub"
	apiKeyCommand      = "/apikey"
	intervalCommand    = "/interval"
)

func NewBot(
	config config.Config,
	server *server.Server,
	notificationCh chan notify.Notification,
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

func (b *Bot) handleNotifications() {
	workerCount := runtime.NumCPU()
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for notification := range b.notificationCh {
				b.SendMessage(notification.ChatID, notification.Message)
			}
		}()
	}

	wg.Wait()
}
