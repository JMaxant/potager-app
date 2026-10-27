// Package web holds the HTTP handlers and middleware.
package web

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Store is what the handlers need from the database, declared here where it is used.
type Store interface {
	Ping(ctx context.Context) error
}

// Server holds the dependencies shared by all handlers.
type Server struct {
	store Store
	log   *slog.Logger
}

// NewServer builds a Server from its logger and its store.
func NewServer(log *slog.Logger, store Store) *Server {
	return &Server{
		store: store,
		log:   log,
	}
}

// Handler returns the router, wrapped in the request-logging middleware.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)

	return logRequests(s.log)(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Pings postgres db, fails if unavailable.
	if err := s.store.Ping(ctx); err != nil {
		s.log.Error("database unavailable", "err", err)
		http.Error(w, "health check failed", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok"))
}
