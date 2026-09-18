package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging_PassesResponseThroughUnchanged(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("hello"))
	})

	handler := Logging(logger)(next)

	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusTeapot)
	}
	if rec.Body.String() != "hello" {
		t.Errorf("body = %q, want %q", rec.Body.String(), "hello")
	}
}

func TestLogging_LogsRequestDetails(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	handler := Logging(logger)(next)

	req := httptest.NewRequest(http.MethodPost, "/api/hello", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log output as JSON: %v\noutput: %s", err, buf.String())
	}

	if entry["method"] != http.MethodPost {
		t.Errorf("method = %v, want %v", entry["method"], http.MethodPost)
	}
	if entry["path"] != "/api/hello" {
		t.Errorf("path = %v, want %v", entry["path"], "/api/hello")
	}
	if entry["status"] != float64(http.StatusCreated) {
		t.Errorf("status = %v, want %v", entry["status"], http.StatusCreated)
	}
	if _, ok := entry["duration_ms"]; !ok {
		t.Error("expected a duration_ms field in the log entry")
	}
}

func TestLogging_DefaultsStatusToOKWhenNotSet(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	})

	handler := Logging(logger)(next)

	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log output as JSON: %v", err)
	}

	if entry["status"] != float64(http.StatusOK) {
		t.Errorf("status = %v, want %v", entry["status"], http.StatusOK)
	}
}

func TestLogging_SkipsHealthEndpoint(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Logging(logger)(next)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	if buf.Len() != 0 {
		t.Errorf("expected no log output for /health, got: %s", buf.String())
	}
}
