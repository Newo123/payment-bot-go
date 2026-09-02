package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Newo123/payment-bot-go/internal/features/payments"
	"github.com/Newo123/payment-bot-go/internal/features/users"
	"github.com/Newo123/payment-bot-go/internal/infrastructure/logger"
	"github.com/Newo123/payment-bot-go/internal/infrastructure/logger/zap"
	"github.com/Newo123/payment-bot-go/internal/infrastructure/postgres/pgx"
	"github.com/Newo123/payment-bot-go/internal/infrastructure/redis/goredis"
	"github.com/Newo123/payment-bot-go/internal/transport/telegram"
)

func main() {
	time.Local = time.UTC

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer cancel()

	log, err := zap.NewLogger(zap.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	log.Debug("application time zone", logger.Any("zone", time.Local))

	postgresPool, err := pgx.NewPool(ctx, pgx.NewConfigMust())
	if err != nil {
		log.Fatal("failed to init postgres connection pool", logger.Error(err))
	}
	defer postgresPool.Close()
	log.Debug("initializing postgres connection pool")

	redisPool, err := goredis.NewPool(ctx, goredis.NewConfigMust())
	if err != nil {
		log.Fatal("failed to init redis connection pool", logger.Error(err))
	}
	defer redisPool.Close()
	log.Debug("initializing redis connection pool")

	// ----- Dependency Injection Start -----

	// Repositories
	usersPGRepo := users.NewPostgresRepository(postgresPool)

	// UseCases
	usersUC := users.NewUseCase(usersPGRepo)

	// TG Handlers
	usersHandlerTG := users.NewHandlerTG(usersUC)
	paymentsHandlerTG := payments.NewHandlerTG()
	// paymentsHandler := paymentstelegram.NewHandler()

	// ----- Dependency Injection End -----

	telegramBot, err := telegram.NewBot(telegram.NewConfigMust(), log)
	if err != nil {
		log.Fatal("failed to init Telegram bot", logger.Error(err))
	}

	// ----- Add Routes Start -----
	telegramBot.RegisterRoutes(
		usersHandlerTG.Routes()...,
	)
	telegramBot.RegisterRoutes(
		paymentsHandlerTG.Routes()...,
	)
	// ----- Add Routes End -----

	if err := telegramBot.Run(ctx); err != nil {
		log.Error("Telegram bot run error", logger.Error(err))
	}
}
