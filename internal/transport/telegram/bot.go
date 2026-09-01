package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Newo123/payment-bot-go/internal/infrastructure/logger"
)

type Bot struct {
	client     *Client
	log        logger.Logger
	config     Config
	middleware []Middleware
	routes     []registeredRoute
}

type registeredRoute struct {
	route   Route
	handler Handler
}

func NewBot(
	config Config,
	log logger.Logger,
	middleware ...Middleware,
) *Bot {
	return &Bot{
		client:     NewClient(config),
		log:        log,
		config:     config,
		middleware: middleware,
	}
}

func (b *Bot) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		if route.Handler == nil {
			b.log.Warn(
				"[TGBOT] skip route without handler",
				logger.String(
					"pattern",
					route.Pattern,
				),
			)

			continue
		}

		handler := chain(
			route.Handler,
			append(
				b.middleware,
				route.Middleware...,
			)...,
		)

		b.routes = append(
			b.routes,
			registeredRoute{
				route: route,

				handler: handler,
			},
		)

		b.log.Info(
			"[TGBOT] route registered",
			logger.String(
				"pattern",
				route.Pattern,
			),
		)
	}
}

func (b *Bot) SendMessageText(ctx context.Context, chatID int64, text string) error {
	_, err := b.client.SendMessage(
		ctx,
		SendMessageParams{
			ChatID: chatID,
			Text:   text,
		},
	)

	if err != nil {
		return fmt.Errorf(
			"send message text: %w",
			err,
		)
	}

	return nil
}

func (b *Bot) SendMessage(ctx context.Context, params SendMessageParams) (Message, error) {
	message, err := b.client.SendMessage(
		ctx,
		params,
	)
	if err != nil {
		return Message{}, fmt.Errorf(
			"send message: %w",
			err,
		)
	}

	return message, nil
}

func (b *Bot) EditMessageText(
	ctx context.Context,
	chatID int64,
	messageID int64,
	text string,
	keyboard *InlineKeyboardMarkup,
) error {
	_, err := b.client.EditMessageText(
		ctx,
		EditMessageTextParams{
			ChatID: chatID,

			MessageID: messageID,

			Text: text,

			ReplyMarkup: keyboard,
		},
	)

	if err != nil {
		return fmt.Errorf(
			"edit message text: %w",
			err,
		)
	}

	return nil
}

func (b *Bot) DeleteMessage(ctx context.Context, chatID int64, messageID int64) error {
	if err := b.client.DeleteMessage(
		ctx,
		DeleteMessageParams{
			ChatID: chatID,

			MessageID: messageID,
		},
	); err != nil {
		return fmt.Errorf(
			"delete message: %w",
			err,
		)
	}

	return nil
}

func (b *Bot) AnswerCallbackQuery(ctx context.Context, callbackQueryID string) error {
	if err := b.client.AnswerCallbackQuery(
		ctx,
		AnswerCallbackQueryParams{
			CallbackQueryID: callbackQueryID,
		},
	); err != nil {
		return fmt.Errorf(
			"answer callback query: %w",
			err,
		)
	}

	return nil
}

func (b *Bot) dispatch(ctx context.Context, update *Update) {
	if update == nil {
		return
	}

	for _, route := range b.routes {
		if !match(
			route.route,
			update,
		) {
			continue
		}

		t := Context{
			Context: ctx,

			Update: update,

			Bot: b,
		}

		if err := route.handler(t); err != nil {
			b.log.Error(
				"[TGBOT] handler failed",
				logger.Error(err),
			)
		}

		return
	}
}

func match(route Route, update *Update) bool {
	switch route.MatchType {
	case MatchExact,
		MatchPrefix,
		MatchContains,
		MatchCommand,
		MatchCommandStartOnly:

		return matchMessage(
			route,
			update,
		)

	case MatchCallbackData:

		return matchCallback(
			route,
			update,
		)

	default:
		return false
	}
}

func matchMessage(route Route, update *Update) bool {
	if update == nil ||
		update.Message == nil {
		return false
	}

	text := strings.TrimSpace(
		update.Message.Text,
	)

	switch route.MatchType {
	case MatchExact:
		return text == route.Pattern

	case MatchPrefix:
		return strings.HasPrefix(
			text,
			route.Pattern,
		)

	case MatchContains:
		return strings.Contains(
			text,
			route.Pattern,
		)

	case MatchCommand:
		return strings.HasPrefix(
			text,
			route.Pattern,
		)

	case MatchCommandStartOnly:
		return isStartCommand(
			text,
			route.Pattern,
		)

	default:
		return false
	}
}

func matchCallback(route Route, update *Update) bool {
	if update == nil ||
		update.CallbackQuery == nil {
		return false
	}

	return update.CallbackQuery.Data ==
		route.Pattern
}

func isStartCommand(text string, pattern string) bool {
	if text == pattern {
		return true
	}

	return strings.HasPrefix(
		text,
		pattern+"@",
	)
}

func (b *Bot) waitRetry(ctx context.Context) bool {
	timer := time.NewTimer(
		b.config.RetryDelay,
	)

	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false

	case <-timer.C:
		return true
	}
}
