package telegram

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token           string        `envconfig:"BOT_TOKEN" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("TELEGRAM", &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}
