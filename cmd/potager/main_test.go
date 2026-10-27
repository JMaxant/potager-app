package main

import (
	"log/slog"
	"strings"
	"testing"
)

const dbURL = "postgres://db:db@127.0.0.1:5432/db"

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{
			name: "all variables set",
			env: map[string]string{
				"DATABASE_URL": dbURL,
				"LISTEN_ADDR":  ":8081",
				"LOG_LEVEL":    "debug",
			},
			want: Config{
				DatabaseURL: dbURL,
				ListenAddr:  ":8081",
				LogLevel:    slog.LevelDebug,
			},
		},
		{
			name: "missing log level, defaults to info",
			env: map[string]string{
				"DATABASE_URL": dbURL,
				"LISTEN_ADDR":  ":8081",
			},
			want: Config{
				DatabaseURL: dbURL,
				ListenAddr:  ":8081",
				LogLevel:    slog.LevelInfo,
			},
		},
		{
			name: "log level warn",
			env: map[string]string{
				"DATABASE_URL": dbURL,
				"LISTEN_ADDR":  ":8081",
				"LOG_LEVEL":    "warn",
			},
			want: Config{
				DatabaseURL: dbURL,
				ListenAddr:  ":8081",
				LogLevel:    slog.LevelWarn,
			},
		},
		{
			name: "log level error",
			env: map[string]string{
				"DATABASE_URL": dbURL,
				"LISTEN_ADDR":  ":8081",
				"LOG_LEVEL":    "error",
			},
			want: Config{
				DatabaseURL: dbURL,
				ListenAddr:  ":8081",
				LogLevel:    slog.LevelError,
			},
		},
		{
			name: "log level uppercase",
			env: map[string]string{
				"DATABASE_URL": dbURL,
				"LISTEN_ADDR":  ":8081",
				"LOG_LEVEL":    "INFO",
			},
			want: Config{
				DatabaseURL: dbURL,
				ListenAddr:  ":8081",
				LogLevel:    slog.LevelInfo,
			},
		},
		{
			name: "log level capitalized",
			env: map[string]string{
				"DATABASE_URL": dbURL,
				"LISTEN_ADDR":  ":8081",
				"LOG_LEVEL":    "Debug",
			},
			want: Config{
				DatabaseURL: dbURL,
				ListenAddr:  ":8081",
				LogLevel:    slog.LevelDebug,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string { return tt.env[key] }
			cfg, err := loadConfig(getenv)

			if err != nil {
				t.Fatal(err)
			}
			if cfg != tt.want {
				t.Errorf("got %+v, want %+v", cfg, tt.want)
			}
		})
	}
}

func TestLoadConfigErrors(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr []string
	}{
		{
			name: "missing database url",
			env: map[string]string{
				"LISTEN_ADDR": ":8081",
			},
			wantErr: []string{"DATABASE_URL"},
		},
		{
			name: "missing listen addr",
			env: map[string]string{
				"DATABASE_URL": dbURL,
			},
			wantErr: []string{"LISTEN_ADDR"},
		},
		{
			name: "invalid log level",
			env: map[string]string{
				"DATABASE_URL": dbURL,
				"LISTEN_ADDR":  ":8081",
				"LOG_LEVEL":    "invalid",
			},
			wantErr: []string{"LOG_LEVEL"},
		},
		{
			name: "all errors",
			env: map[string]string{
				"LOG_LEVEL": "invalid",
			},
			wantErr: []string{"DATABASE_URL", "LISTEN_ADDR", "LOG_LEVEL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string {
				return tt.env[key]
			}
			cfg, err := loadConfig(getenv)
			if cfg != (Config{}) {
				t.Errorf("got %+v, want empty config", cfg)
			}
			if err == nil {
				t.Fatal("expected error")
			}

			joined, ok := err.(interface{ Unwrap() []error })
			if !ok {
				t.Fatalf("error %T does not wrap several errors", err)
			}
			if got := len(joined.Unwrap()); got != len(tt.wantErr) {
				t.Errorf("got %d errors, want %d: %v", got, len(tt.wantErr), err)
			}

			for _, substr := range tt.wantErr {
				if !strings.Contains(err.Error(), substr) {
					t.Errorf("got error %q, want to contain %q", err.Error(), substr)
				}
			}
		})
	}
}
