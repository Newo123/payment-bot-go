package goredis

import (
	"errors"

	"github.com/Newo123/payment-bot-go/internal/infrastructure/redis"
	goredis "github.com/redis/go-redis/v9"
)

type goredisStringCmd struct {
	*goredis.StringCmd
}

func (c goredisStringCmd) Bytes() ([]byte, error) {
	data, err := c.StringCmd.Bytes()
	if err != nil {
		return nil, mapErrors(err)
	}

	return data, nil
}

type goredisStatusCmd struct {
	*goredis.StatusCmd
}

type goredisIntCmd struct {
	*goredis.IntCmd
}

func mapErrors(err error) error {
	if errors.Is(err, redis.NotFound) {
		return redis.NotFound
	}

	return err
}
