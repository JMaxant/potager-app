package web

import (
	"log/slog"
	"net/http"
)

func logRequests(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Instantiate StatusRecorder with default values.
			writer := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}
			// Serve the request with the StatusRecorder instead of the ResponseWriter.
			next.ServeHTTP(writer, r)
			// During the request, the StatusRecorder will capture the status code, overriding the default value.
			logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", writer.status)
		})
	}
}
