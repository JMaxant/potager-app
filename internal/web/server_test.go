package web

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeStore struct {
	err error
}

func (f fakeStore) Ping(ctx context.Context) error {
	return f.err
}

func TestServer_handleHealth(t *testing.T) {
	tests := []struct {
		name     string
		storeErr error
		wantCode int
		wantBody string
	}{
		{
			name:     "database up",
			storeErr: nil,
			wantCode: http.StatusOK,
			wantBody: "ok",
		},
		{
			name:     "database down",
			storeErr: errors.New("database is down"),
			wantCode: http.StatusServiceUnavailable,
			wantBody: "health check failed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			srv := NewServer(logger, fakeStore{err: tt.storeErr})

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()
			srv.handleHealth(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("got status %d, want %d", w.Code, tt.wantCode)
			}
			if w.Body.String() != tt.wantBody {
				t.Errorf("got body %q, want %q", w.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestServer_Handler(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		wantCode int
	}{
		{
			name:     "GET /health",
			method:   http.MethodGet,
			path:     "/health",
			wantCode: http.StatusOK,
		},
		{
			name:     "POST /health",
			method:   http.MethodPost,
			path:     "/health",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "GET /nope",
			method:   http.MethodGet,
			path:     "/nope",
			wantCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			srv := NewServer(logger, fakeStore{})

			req := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()
			srv.Handler().ServeHTTP(w, req)
			if w.Code != test.wantCode {
				t.Errorf("got status %d, want %d", w.Code, test.wantCode)
			}
		})
	}
}
