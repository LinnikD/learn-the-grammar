// Package middleware provides HTTP middleware shared across the backend.
package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/LinnikD/learn-the-grammar/backend/internal/apierror"
)

// skipLogging holds paths excluded from request logging, namely the
// health check endpoint that Kubernetes readiness/liveness probes hit
// every few seconds and would otherwise flood the logs.
var skipLogging = map[string]bool{
	"/health": true,
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Logging returns middleware that logs each request's method, path,
// status code, and duration using logger.
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if skipLogging[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()

			next.ServeHTTP(rec, r)

			logger.Info("http request",
				"request_id", apierror.RequestID(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", float64(time.Since(start).Microseconds())/1000,
			)
		})
	}
}
