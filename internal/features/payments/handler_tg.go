package payments

import (
	"github.com/Newo123/payment-bot-go/internal/transport/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handlerTG struct {
}

func NewHandlerTG() *handlerTG {
	return &handlerTG{}
}

func (h *handlerTG) Routes() []telegram.Route {
	return []telegram.Route{
		{
			Match:   telegram.Text("📥 Пополнить счет"),
			Handler: h.Deposit,
		},
		{
			Match:   telegram.Text("📤 Вывести средства"),
			Handler: h.Payout,
		},
	}
}

func (h *handlerTG) Deposit(ctx telegram.Context) error {
	return ctx.Send(tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "DEPOSIT"))
}
func (h *handlerTG) Payout(ctx telegram.Context) error {
	return ctx.Send(tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "PAYOUT"))
}
