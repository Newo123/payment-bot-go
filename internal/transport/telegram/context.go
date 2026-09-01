package telegram

import "context"

type Context struct {
	context.Context

	Update *Update

	Bot *Bot
}
