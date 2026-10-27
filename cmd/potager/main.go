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

// Config is the process configuration, read once from the environment at start-up.
type Config struct {
	DatabaseURL string
	ListenAddr  string
	LogLevel    slog.Level
}

func main() {
	// main only turns run's error into an exit code, so run can use defer (os.Exit skips defers).
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// config initialization, if any errors are met, abort.
	config, err := loadConfig(os.Getenv)
	if err != nil {
		return fmt.Errorf("config:\n %w", err)
	}

	// ctx is cancelled on Ctrl+C or SIGTERM (docker stop), which triggers the shutdown below.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// the database connection gets 10 seconds, then the attempt is abandoned.
	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	db, err := store.New(connectCtx, config.DatabaseURL)
	// the timeout only matters during the connection: release its timer right away.
	cancel()

	if err != nil {
		return fmt.Errorf("database: %w", err)
	}

	// close the connection pool when run returns.
	defer db.Close()

	// structured JSON logs on stdout, filtered by the configured level.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.LogLevel}))
	logger.Info("starting server", "listen_addr", config.ListenAddr)

	// wire the dependencies by hand: the web layer receives its logger and its database.
	server := web.NewServer(logger, db)
	// timeouts protect against slow or idle clients holding connections open.
	srv := &http.Server{
		Addr:              config.ListenAddr,
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	// the server runs in a goroutine; this channel carries its failure back to run.
	errs := make(chan error, 1)

	go func() {
		logger.Info("listening", "listen_addr", config.ListenAddr)
		// ErrServerClosed is the normal result of Shutdown, not a failure.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- fmt.Errorf("serve: %w", err)
		}
	}()

	// block here until the server fails or a stop signal arrives, whichever comes first.
	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		logger.Info("shutting down")
	}

	// fresh context: ctx is already cancelled, and in-flight requests get 10 seconds to finish.
	webCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	// stop accepting connections, then wait for the running requests.
	err = srv.Shutdown(webCtx)

	if err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}

// loadConfig takes getenv as a parameter so tests can pass a fake instead of the real environment.
func loadConfig(getenv func(string) string) (Config, error) {
	databaseURL := getenv("DATABASE_URL")
	// collect every problem instead of stopping at the first one.
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

	// errors.Join merges all the problems into one error, one per line.
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

	// unset means the zero value of slog.Level, which is Info.
	if level == "" {
		return lvl, nil
	}
	// accepts "debug", "info", "warn" and "error" (case-insensitive).
	err := lvl.UnmarshalText([]byte(level))
	return lvl, err
}
