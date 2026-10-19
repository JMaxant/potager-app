package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
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
	slog.SetDefault(logger)

	logger.Info("starting server", "listen_addr", config.ListenAddr)

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

	config := Config{
		DatabaseURL: databaseURL,
		ListenAddr:  listenAddr,
		LogLevel:    level,
	}

	return config, errors.Join(errs...)
}

func parseLogLevel(level string) (slog.Level, error) {
	var lvl slog.Level

	if level == "" {
		return lvl, nil
	}
	err := lvl.UnmarshalText([]byte(level))
	return lvl, err
}
