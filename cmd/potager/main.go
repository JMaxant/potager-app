package main

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

	"github.com/JMaxant/potager-app/internal/store"
	"github.com/JMaxant/potager-app/internal/web"
)

type Config struct {
	DatabaseURL string
	ListenAddr  string
	LogLevel    slog.Level
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := loadConfig(os.Getenv)
	if err != nil {
		return fmt.Errorf("config:\n %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.LogLevel}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	db, err := store.New(connectCtx, config.DatabaseURL)
	cancel()

	if err != nil {
		return fmt.Errorf("database: %w", err)
	}

	defer db.Close()

	logger.Info("starting server", "listen_addr", config.ListenAddr)

	server := web.NewServer(logger, db)
	srv := &http.Server{
		Addr:              config.ListenAddr,
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	errs := make(chan error, 1)

	go func() {
		logger.Info("listening", "listen_addr", config.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- fmt.Errorf("serve: %w", err)
		}
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		logger.Info("shutting down")
	}

	webCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = srv.Shutdown(webCtx)

	if err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}

func loadConfig(getenv func(string) string) (Config, error) {
	databaseURL := getenv("DATABASE_URL")
	var errs []error

	if databaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}
	listenAddr := getenv("LISTEN_ADDR")
	if listenAddr == "" {
		errs = append(errs, errors.New("LISTEN_ADDR is required"))
	}

	level, err := parseLogLevel(getenv("LOG_LEVEL"))

	if err != nil {
		errs = append(errs, fmt.Errorf("LOG_LEVEL: %w", err))
	}

	if len(errs) > 0 {
		return Config{}, errors.Join(errs...)
	}

	config := Config{
		DatabaseURL: databaseURL,
		ListenAddr:  listenAddr,
		LogLevel:    level,
	}

	return config, nil
}

func parseLogLevel(level string) (slog.Level, error) {
	var lvl slog.Level

	if level == "" {
		return lvl, nil
	}
	err := lvl.UnmarshalText([]byte(level))
	return lvl, err
}
