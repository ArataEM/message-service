package application

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ArataEM/message-service/config"
	"github.com/ArataEM/message-service/repository"
)

type App struct {
	router http.Handler
	rdb    repository.Repo
	config config.Config
}

func New(config config.Config) *App {
	app := &App{
		rdb:    repository.NewRedisRepo(config),
		config: config,
	}

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.ServerPort),
		Handler:      a.router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	err := a.rdb.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to repository: %w", err)
	}

	defer func() {
		err := a.rdb.Close()
		if err != nil {
			fmt.Println("failed to close connection: ", err)
		}
	}()

	slog.Info("Starting server")

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("server error: %w", err)
		}
		close(ch)
	}()

	slog.Info(fmt.Sprintf("Server started on port %d", a.config.ServerPort))

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		slog.Info("Graceful shutdown")
		return server.Shutdown(timeout)
	}
}
