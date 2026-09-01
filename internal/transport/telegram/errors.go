package telegram

import "fmt"

// APIError — ошибка, которую вернул Telegram Bot API.
type APIError struct {
	Code        int
	Description string
}

func (e *APIError) Error() string {
	return fmt.Sprintf(
		"telegram api error: code=%d description=%s",
		e.Code,
		e.Description,
	)
}
