// Package web holds the HTTP surface: routing, middleware and, from step 1
// onwards, the templ components. It depends on the domain, never the reverse.
package web

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Store is the database surface the web layer needs. It is declared here, where
// it is consumed, rather than in the store package that implements it — the Go
// convention, and the one most Java-shaped Go codebases get backwards.
type Store interface {
	Ping(ctx context.Context) error
}

type Server struct {
	log   *slog.Logger
	store Store
	mux   *http.ServeMux
}

func NewServer(log *slog.Logger, store Store) *Server {
	s := &Server{log: log, store: store, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.requestLogger(s.mux).ServeHTTP(w, r)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if err := s.store.Ping(r.Context()); err != nil {
		s.log.ErrorContext(r.Context(), "health check failed", "error", err)
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		s.log.InfoContext(r.Context(), "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
