package web

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogRequests(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus int
	}{
		{
			name: "explicit 404",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "no explicit status",
			handler: func(w http.ResponseWriter, r *http.Request) {
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "first WriteHeader wins",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				w.WriteHeader(http.StatusOK)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "body without WriteHeader",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hi"))
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

			h := logRequests(logger)(tt.handler)
			req := httptest.NewRequest(http.MethodGet, "/plots", nil)
			h.ServeHTTP(httptest.NewRecorder(), req)

			var got map[string]any
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("log is not valid JSON: %v\n%s", err, buf.String())
			}
			if got["status"] != float64(tt.wantStatus) {
				t.Errorf("got status %v, want %d", got["status"], tt.wantStatus)
			}
			if got["method"] != http.MethodGet {
				t.Errorf("got method %v, want GET", got["method"])
			}

			if _, ok := got["duration"]; !ok {
				t.Error("duration is missing for the log")
			}

			if got["path"] != "/plots" {
				t.Errorf("got path %v, want /plots", got["path"])
			}
		})
	}
}
