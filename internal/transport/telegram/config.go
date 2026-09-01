package telegram

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token string `envconfig:"BOT_TOKEN" required:"true"`

	BaseURL string `envconfig:"BASE_URL" default:"https://api.telegram.org"`

	// HTTPTimeout должен быть больше PollingTimeout.
	//
	// Например:
	// HTTPTimeout     = 60s
	// PollingTimeout  = 30s
	//
	// Иначе http.Client может завершить getUpdates раньше Telegram.
	HTTPTimeout time.Duration `envconfig:"HTTP_TIMEOUT" default:"60s"`

	PollingTimeout time.Duration `envconfig:"POLLING_TIMEOUT" default:"30s"`

	PollingLimit int `envconfig:"POLLING_LIMIT" default:"100"`

	// Если Telegram временно недоступен,
	// бот ждёт это время перед следующей попыткой.
	RetryDelay time.Duration `envconfig:"RETRY_DELAY" default:"1s"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process(
		"TELEGRAM",
		&config,
	); err != nil {
		return Config{}, fmt.Errorf(
			"process Telegram config: %w",
			err,
		)
	}

	if config.PollingTimeout <= 0 {
		config.PollingTimeout = 30 * time.Second
	}

	if config.PollingLimit <= 0 || config.PollingLimit > 100 {
		config.PollingLimit = 100
	}

	// Защищаемся от типичной ошибки:
	// HTTP timeout <= long polling timeout.
	if config.HTTPTimeout <= config.PollingTimeout {
		config.HTTPTimeout = config.PollingTimeout + 30*time.Second
	}

	if config.RetryDelay <= 0 {
		config.RetryDelay = time.Second
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(
			fmt.Errorf(
				"get Telegram config: %w",
				err,
			),
		)
	}

	return config
}
