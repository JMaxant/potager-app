package main

import (
	"log/slog"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	const dbURL = "postgres://db:db@127.0.0.1:5432/db"
	tests := []struct {
		name    string
		env     map[string]string
		want    Config
		wantErr []string
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
			getenv := func(key string) string { return tt.env[key] }
			cfg, err := loadConfig(getenv)

			if len(tt.wantErr) == 0 {
				if err != nil {
					t.Fatal(err)
				}
				if cfg != tt.want {
					t.Errorf("got %+v, want %+v", cfg, tt.want)
				}
			} else {
				if err == nil {
					t.Fatal("expected error")
				}
				for _, substr := range tt.wantErr {
					if !strings.Contains(err.Error(), substr) {
						t.Errorf("got error %q, want to contain %q", err.Error(), substr)
					}
				}
			}
		})
	}
}
