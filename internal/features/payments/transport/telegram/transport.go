package telegram

import (
	apptelegram "github.com/Newo123/payment-bot-go/internal/transport/telegram"
)

const (
	paymentCreateCallback = "payment:create"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Routes() []apptelegram.Route {
	return []apptelegram.Route{
		{
			Pattern:   "/pay",
			MatchType: apptelegram.MatchCommandStartOnly,
			Handler:   h.Pay,
		},
		{
			Pattern:   paymentCreateCallback,
			MatchType: apptelegram.MatchCallbackData,
			Handler:   h.CreatePayment,
		},
	}
}

func (h *Handler) Pay(
	t apptelegram.Context,
) error {
	if t.Update == nil || t.Update.Message == nil {
		return nil
	}

	_, err := t.Bot.SendMessage(
		t.Context,
		apptelegram.SendMessageParams{
			ChatID: t.Update.Message.Chat.ID,
			Text:   "Выберите действие:",
			ReplyMarkup: &apptelegram.InlineKeyboardMarkup{
				InlineKeyboard: [][]apptelegram.InlineKeyboardButton{
					{
						{
							Text:         "💳 Оплатить",
							CallbackData: paymentCreateCallback,
						},
					},
					{
						{
							Text: "📱 Открыть приложение",
							WebApp: &apptelegram.WebAppInfo{
								URL: "https://example.com/app",
							},
						},
					},
				},
			},
		},
	)

	return err
}

func (h *Handler) CreatePayment(
	t apptelegram.Context,
) error {
	if t.Update == nil || t.Update.CallbackQuery == nil {
		return nil
	}

	callback := t.Update.CallbackQuery

	// Убираем у пользователя индикатор загрузки
	// после нажатия inline-кнопки.
	if err := t.Bot.AnswerCallbackQuery(
		t.Context,
		callback.ID,
	); err != nil {
		return err
	}

	if callback.Message == nil {
		return nil
	}

	return t.Bot.EditMessageText(
		t.Context,
		callback.Message.Chat.ID,
		callback.Message.ID,
		"✅ Платёж создан",
		nil,
	)
}
