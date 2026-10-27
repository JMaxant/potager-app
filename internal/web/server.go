package web

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Store interface {
	Ping(ctx context.Context) error
}

type Server struct {
	store Store
	log   *slog.Logger
}

func NewServer(log *slog.Logger, store Store) *Server {
	return &Server{
		store: store,
		log:   log,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)

	return logRequests(s.log)(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.store.Ping(ctx); err != nil {
		s.log.Error("database unavailable", "err", err)
		http.Error(w, "health check failed", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok"))
}
