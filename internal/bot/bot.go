package bot

import (
	"runtime"
	"sync"

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

func (b *Bot) handleNotifications() {
	workerCount := runtime.NumCPU()

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for notification := range b.notificationCh {
				if notification.Message != "" {
					if err := b.SendMessage(
						notification.ChatID,
						notification.Message,
					); err != nil {
						for _, update := range notification.Updates {
							updateStorage(b.store, notification.ChatID, update)
						}
					}
				} else if len(notification.Updates) > 0 {
					for _, update := range notification.Updates {
						updateStorage(b.store, notification.ChatID, update)
					}
				}
			}
		}()
	}

	wg.Wait()
}

func updateStorage(
	store *storage.Storage,
	chatID int64,
	update notification.Update,
) {
	if update.LatestTimestamp != 0 {
		if err := store.SetLatestSubtitleTimestamp(
			chatID,
			update.TitleID,
			update.LatestTimestamp,
		); err != nil {
			log.Error(
				"failed to set latest timestamp",
				"chatID",
				chatID,
				"update",
				update,
				"err",
				err)
		}
	}

	if update.JapaneseName != "" {
		if err := store.SetJapaneseName(
			chatID,
			update.TitleID,
			update.JapaneseName,
		); err != nil {
			log.Error(
				"failed to set japanese name",
				"chatID",
				chatID,
				"update",
				update,
				"err",
				err)
		}
	}
}
