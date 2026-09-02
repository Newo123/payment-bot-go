package users

import (
	"errors"
	"fmt"

	"github.com/Newo123/payment-bot-go/internal/domain"
	"github.com/Newo123/payment-bot-go/internal/transport/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handlerTG struct {
	useCase UseCase
}

func NewHandlerTG(useCase UseCase) *handlerTG {
	return &handlerTG{
		useCase: useCase,
	}
}

func (h *handlerTG) Routes() []telegram.Route {
	return []telegram.Route{
		{
			Match:   telegram.Command("start"),
			Handler: h.Start,
		},
		{
			Match:   telegram.Contact(),
			Handler: h.Contact,
		},
	}
}

func (h *handlerTG) Start(ctx telegram.Context) error {
	_, err := h.useCase.FindByID(ctx.Context, ctx.Update.Message.From.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return h.requestContact(ctx)
		}

		return fmt.Errorf("find user by telegram id: %w", err)
	}

	_, err = h.useCase.Update(
		ctx.Context,
		ctx.Update.Message.From.ID,
		UpdateParams{
			FirstName:    ctx.Update.Message.From.FirstName,
			LastName:     ctx.Update.Message.From.LastName,
			Username:     ctx.Update.Message.From.UserName,
			LanguageCode: ctx.Update.Message.From.LanguageCode,
		},
	)
	if err != nil {
		return fmt.Errorf("update telegram user: %w", err)
	}

	return h.showMainMenu(ctx)
}

func (h *handlerTG) Contact(ctx telegram.Context) error {
	_, err := h.useCase.Create(ctx.Context, CreateParams{
		ID:           ctx.Update.Message.From.ID,
		Phone:        ctx.Update.Message.Contact.PhoneNumber,
		FirstName:    ctx.Update.Message.From.FirstName,
		LastName:     ctx.Update.Message.From.LastName,
		LanguageCode: ctx.Update.Message.From.LanguageCode,
		Username:     ctx.Update.Message.From.UserName,
	})
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return h.showMainMenu(ctx)
}

func (h *handlerTG) requestContact(ctx telegram.Context) error {
	msg := tgbotapi.NewMessage(
		ctx.Update.Message.Chat.ID,
		"Для начала работы отправьте свой номер телефона! 👇",
	)

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("📞 Отправить номер телефона"),
		),
	)

	return ctx.Send(msg)
}

func (h *handlerTG) showMainMenu(ctx telegram.Context) error {
	msg := tgbotapi.NewMessage(ctx.Update.Message.From.ID, `
	⚡️ Быстрое и удобное пополнение и вывод средств
	
	— Обработка без задержек
	— Комиссия 0%
	
	Выберите раздел ниже 👇🏻
	`)

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📥 Пополнить счет"),
			tgbotapi.NewKeyboardButton("📤 Вывести средства"),
		),
	)

	return ctx.Send(msg)
}
