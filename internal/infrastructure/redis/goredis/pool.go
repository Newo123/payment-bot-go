package goredis

import (
	"context"
	"fmt"
	"time"

	"github.com/Newo123/payment-bot-go/internal/infrastructure/redis"
	goredis "github.com/redis/go-redis/v9"
)

type Pool struct {
	client *goredis.Client
	ttl    time.Duration
}

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	options := &goredis.Options{
		Addr:     config.Addr(),
		Password: config.Password,
		DB:       config.Database,
	}

	client := goredis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Pool{
		client: client,
		ttl:    config.TTL,
	}, nil
}

func (p *Pool) Get(ctx context.Context, key string) redis.StatusCmd {
	cmd := p.client.Get(ctx, key)

	return goredisStringCmd{cmd}
}

func (p *Pool) Set(ctx context.Context, key string, value any, ttl time.Duration) redis.StatusCmd {
	cmd := p.client.Set(ctx, key, value, ttl)

	return goredisStatusCmd{cmd}
}

func (p *Pool) Del(ctx context.Context, keys ...string) redis.IntCmd {
	cmd := p.client.Del(ctx, keys...)

	return goredisIntCmd{cmd}
}

func (p *Pool) HGet(ctx context.Context, key string, field string) redis.StringCmd {
	cmd := p.client.HGet(ctx, key, field)

	return goredisStringCmd{cmd}
}

func (p *Pool) HSet(ctx context.Context, key string, values ...any) redis.IntCmd {
	cmd := p.client.HSet(ctx, key, values...)

	return goredisIntCmd{cmd}
}

func (p *Pool) Close() error {
	return p.client.Close()
}

func (p *Pool) TTL() time.Duration {
	return p.ttl
}
