package telegram

func Command(name string) func(Context) bool {
	return func(ctx Context) bool {
		msg := ctx.Update.Message

		return msg != nil &&
			msg.IsCommand() &&
			msg.Command() == name
	}
}

func Contact() func(Context) bool {
	return func(ctx Context) bool {
		return ctx.Update.Message != nil &&
			ctx.Update.Message.Contact != nil
	}
}

func Text(text string) func(Context) bool {
	return func(ctx Context) bool {
		return ctx.Update.Message != nil &&
			ctx.Update.Message.Text == text
	}
}

func Callback(data string) func(Context) bool {
	return func(ctx Context) bool {
		return ctx.Update.CallbackQuery != nil &&
			ctx.Update.CallbackQuery.Data == data
	}
}
