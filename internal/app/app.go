package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IvanKuchsh-600/subtracker/internal/config"
	"github.com/IvanKuchsh-600/subtracker/internal/http/handlers"
	"github.com/IvanKuchsh-600/subtracker/internal/http/router"
	"github.com/IvanKuchsh-600/subtracker/internal/repository/postgres"
	"github.com/IvanKuchsh-600/subtracker/internal/usecase/subscription"
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("create app: %w", errors.New("config is required"))
	}

	return &App{cfg: cfg}, nil
}

func (a *App) Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	connStr := a.cfg.Database.ConnString()
	subscrRepo, err := postgres.NewSubscriptionRepository(connStr, logger)
	if err != nil {
		return err
	}
	logger.Info("Connected to database", "host", a.cfg.Database.Host, "dbname", a.cfg.Database.DBName)

	subscrUsecase := subscription.NewService(subscrRepo, logger)
	subscrHandler := handlers.NewSubscriptionHandler(subscrUsecase)
	router := router.NewRouter(subscrHandler)

	srv := &http.Server{
		Addr:        ":" + a.cfg.Server.Port,
		Handler:     router,
		ReadTimeout: 15 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown http server", "error", err)
		}
	}()

	// logger.Info("http server started", "addr", cfg.HTTPAddr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("listen and serve", "error", err)
		os.Exit(1)
	}
	return nil
}
