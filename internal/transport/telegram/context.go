package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Context struct {
	context.Context
	Bot    *tgbotapi.BotAPI
	Update *tgbotapi.Update
}

func NewContext(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) Context {
	return Context{
		Context: ctx,
		Bot:     bot,
		Update:  update,
	}
}

func (c Context) Send(msg tgbotapi.Chattable) error {
	if _, err := c.Bot.Send(msg); err != nil {
		return fmt.Errorf("send telegram message: %w", err)
	}

	return nil
}

func (c Context) SendText(text string) error {
	if c.Update.Message == nil {
		return nil
	}

	return c.Send(
		tgbotapi.NewMessage(
			c.Update.Message.Chat.ID,
			text,
		),
	)
}
