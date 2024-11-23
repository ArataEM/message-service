package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArataEM/message-service/application"
	"github.com/ArataEM/message-service/config"
)

func main() {
	app := application.New(config.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		slog.Error(err.Error())
	}
}
