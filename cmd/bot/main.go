package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	paymentstelegram "github.com/Newo123/payment-bot-go/internal/features/payments/transport/telegram"
	userstelegram "github.com/Newo123/payment-bot-go/internal/features/users/transport/telegram"
	"github.com/Newo123/payment-bot-go/internal/infrastructure/logger"
	"github.com/Newo123/payment-bot-go/internal/infrastructure/logger/zap"
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

	telegramBot := telegram.NewBot(
		telegram.NewConfigMust(),
		log,
	)

	usersHandler := userstelegram.NewHandler()
	paymentsHandler := paymentstelegram.NewHandler()

	telegramBot.RegisterRoutes(
		usersHandler.Routes()...,
	)

	telegramBot.RegisterRoutes(
		paymentsHandler.Routes()...,
	)

	if err := telegramBot.Start(ctx); err != nil {
		log.Fatal(
			"telegram bot failed",
			logger.Error(err),
		)
	}
}
