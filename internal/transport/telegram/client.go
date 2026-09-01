package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	token   string
	baseURL string
	http    *http.Client
}

func NewClient(config Config) *Client {
	baseURL := strings.TrimRight(
		config.BaseURL,
		"/",
	)

	return &Client{
		token:   config.Token,
		baseURL: baseURL,
		http: &http.Client{
			Timeout: config.HTTPTimeout,
		},
	}
}

type apiResponse[T any] struct {
	OK          bool   `json:"ok"`
	Result      T      `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (c *Client) call(
	ctx context.Context,
	method string,
	params any,
	result any,
) error {
	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf(
			"marshal Telegram request: %w",
			err,
		)
	}

	url := fmt.Sprintf(
		"%s/bot%s/%s",
		c.baseURL,
		c.token,
		method,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf(
			"create Telegram request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf(
			"Telegram HTTP request: %w",
			err,
		)
	}

	defer resp.Body.Close()

	var response apiResponse[json.RawMessage]

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&response); err != nil {
		return fmt.Errorf(
			"decode Telegram response: %w",
			err,
		)
	}

	if !response.OK {
		return &APIError{
			Code:        response.ErrorCode,
			Description: response.Description,
		}
	}

	if result == nil {
		return nil
	}

	if err := json.Unmarshal(
		response.Result,
		result,
	); err != nil {
		return fmt.Errorf(
			"decode Telegram result: %w",
			err,
		)
	}

	return nil
}

type SendMessageParams struct {
	ChatID int64 `json:"chat_id"`

	Text string `json:"text"`

	ReplyMarkup any `json:"reply_markup,omitempty"`
}

func (c *Client) SendMessage(
	ctx context.Context,
	params SendMessageParams,
) (Message, error) {
	var message Message

	if err := c.call(
		ctx,
		"sendMessage",
		params,
		&message,
	); err != nil {
		return Message{}, fmt.Errorf(
			"sendMessage: %w",
			err,
		)
	}

	return message, nil
}

type EditMessageTextParams struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int64  `json:"message_id"`
	Text      string `json:"text"`

	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (c *Client) EditMessageText(
	ctx context.Context,
	params EditMessageTextParams,
) (Message, error) {
	var message Message

	if err := c.call(
		ctx,
		"editMessageText",
		params,
		&message,
	); err != nil {
		return Message{}, fmt.Errorf(
			"editMessageText: %w",
			err,
		)
	}

	return message, nil
}

type DeleteMessageParams struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int64 `json:"message_id"`
}

func (c *Client) DeleteMessage(
	ctx context.Context,
	params DeleteMessageParams,
) error {
	var result bool

	if err := c.call(
		ctx,
		"deleteMessage",
		params,
		&result,
	); err != nil {
		return fmt.Errorf(
			"deleteMessage: %w",
			err,
		)
	}

	return nil
}

type AnswerCallbackQueryParams struct {
	CallbackQueryID string `json:"callback_query_id"`

	Text string `json:"text,omitempty"`

	ShowAlert bool `json:"show_alert,omitempty"`
}

func (c *Client) AnswerCallbackQuery(
	ctx context.Context,
	params AnswerCallbackQueryParams,
) error {
	var result bool

	if err := c.call(
		ctx,
		"answerCallbackQuery",
		params,
		&result,
	); err != nil {
		return fmt.Errorf(
			"answerCallbackQuery: %w",
			err,
		)
	}

	return nil
}

type GetUpdatesParams struct {
	Offset int64 `json:"offset,omitempty"`

	Limit int `json:"limit,omitempty"`

	// Telegram expects seconds here.
	Timeout int `json:"timeout,omitempty"`
}

func (c *Client) GetUpdates(
	ctx context.Context,
	params GetUpdatesParams,
) ([]Update, error) {
	var updates []Update

	if err := c.call(
		ctx,
		"getUpdates",
		params,
		&updates,
	); err != nil {
		return nil, fmt.Errorf(
			"getUpdates: %w",
			err,
		)
	}

	return updates, nil
}

// DeleteWebhook удаляет существующий webhook.
//
// Это важно, если перед запуском polling у бота был установлен webhook.
// Telegram не позволяет одновременно получать updates через webhook
// и getUpdates.
func (c *Client) DeleteWebhook(
	ctx context.Context,
) error {
	var result bool

	if err := c.call(
		ctx,
		"deleteWebhook",
		struct {
			DropPendingUpdates bool `json:"drop_pending_updates,omitempty"`
		}{},
		&result,
	); err != nil {
		return fmt.Errorf(
			"deleteWebhook: %w",
			err,
		)
	}

	return nil
}
