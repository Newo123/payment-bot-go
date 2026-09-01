package telegram

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text string `json:"text"`

	CallbackData string `json:"callback_data,omitempty"`

	URL string `json:"url,omitempty"`

	WebApp *WebAppInfo `json:"web_app,omitempty"`
}

type WebAppInfo struct {
	URL string `json:"url"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`

	IsPersistent bool `json:"is_persistent,omitempty"`

	ResizeKeyboard bool `json:"resize_keyboard,omitempty"`

	OneTimeKeyboard bool `json:"one_time_keyboard,omitempty"`

	Selective bool `json:"selective,omitempty"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`

	Selective bool `json:"selective,omitempty"`
}
