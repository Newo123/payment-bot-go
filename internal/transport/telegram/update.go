package telegram

type Update struct {
	ID int64 `json:"update_id"`

	Message *Message `json:"message,omitempty"`

	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type Message struct {
	ID int64 `json:"message_id"`

	Text string `json:"text"`

	From *User `json:"from,omitempty"`

	Chat Chat `json:"chat"`
}

type User struct {
	ID int64 `json:"id"`

	IsBot bool `json:"is_bot"`

	Username string `json:"username,omitempty"`

	FirstName string `json:"first_name"`

	LastName string `json:"last_name,omitempty"`

	LanguageCode string `json:"language_code,omitempty"`
}

type Chat struct {
	ID int64 `json:"id"`

	Type string `json:"type"`

	Username string `json:"username,omitempty"`

	FirstName string `json:"first_name,omitempty"`

	LastName string `json:"last_name,omitempty"`

	Title string `json:"title,omitempty"`
}

type CallbackQuery struct {
	ID string `json:"id"`

	From User `json:"from"`

	Data string `json:"data"`

	Message *Message `json:"message,omitempty"`
}
