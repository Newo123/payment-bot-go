package telegram

import (
	apptelegram "github.com/Newo123/payment-bot-go/internal/transport/telegram"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Routes() []apptelegram.Route {
	return []apptelegram.Route{
		{
			Pattern:   "/start",
			MatchType: apptelegram.MatchCommandStartOnly,
			Handler:   h.Start,
		},
	}
}

func (h *Handler) Start(
	t apptelegram.Context,
) error {
	return t.Bot.SendMessageText(
		t,
		t.Update.Message.Chat.ID,
		"Привет! 👋",
	)
}
