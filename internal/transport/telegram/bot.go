package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Newo123/payment-bot-go/internal/infrastructure/logger"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	config Config
	log    logger.Logger
	router *Router
}

func NewBot(config Config, log logger.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	return &Bot{
		api:    api,
		config: config,
		log:    log,
		router: NewRouter(),
	}, nil
}

func (b *Bot) RegisterRoutes(routes ...Route) {
	b.router.Register(routes...)
}

func (b *Bot) Run(ctx context.Context) error {
	updates := b.api.GetUpdatesChan(tgbotapi.NewUpdate(0))

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		for update := range updates {
			if err := b.router.Handle(
				NewContext(ctx, b.api, &update),
			); err != nil {
				b.log.Error(
					"telegram update handling error",
					logger.Error(err),
				)
			}
		}
	}()

	select {
	case <-ch:
		return nil

	case <-ctx.Done():
		b.api.StopReceivingUpdates()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			b.config.ShutdownTimeout,
		)
		defer cancel()

		select {
		case <-ch:
			return nil
		case <-shutdownCtx.Done():
			return fmt.Errorf("telegram bot shutdown timeout")
		}
	}
}
