package telegram

import (
	"context"
	"errors"

	"github.com/Newo123/payment-bot-go/internal/infrastructure/logger"
)

func (b *Bot) Start(
	ctx context.Context,
) error {
	b.log.Info(
		"[TGBOT] starting",
	)

	defer func() {
		b.log.Info(
			"[TGBOT] stopped",
		)
	}()

	// Если ранее был webhook,
	// getUpdates работать не будет.
	if err := b.client.DeleteWebhook(ctx); err != nil {
		b.log.Warn(
			"[TGBOT] failed to delete webhook",
			logger.Error(err),
		)
	}

	var offset int64

	timeoutSeconds := int(
		b.config.PollingTimeout.Seconds(),
	)

	for {
		if err := ctx.Err(); err != nil {
			return nil
		}

		updates, err := b.client.GetUpdates(
			ctx,
			GetUpdatesParams{
				Offset: offset,

				Limit: b.config.PollingLimit,

				Timeout: timeoutSeconds,
			},
		)

		if err != nil {
			if errors.Is(
				err,
				context.Canceled,
			) {
				return nil
			}

			if errors.Is(
				err,
				context.DeadlineExceeded,
			) {
				return nil
			}

			b.log.Error(
				"[TGBOT] polling failed",
				logger.Error(err),
			)

			if !b.waitRetry(ctx) {
				return nil
			}

			continue
		}

		for i := range updates {
			update := &updates[i]

			// Telegram рекомендует подтверждать update
			// через update_id + 1.
			offset = update.ID + 1

			go b.dispatch(
				ctx,
				update,
			)
		}
	}
}
