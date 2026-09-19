// Package apierror writes the public API error contract.
package apierror

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LinnikD/learn-the-grammar/backend/internal/api"
)

type requestIDKey struct{}

// WithRequestID attaches the server-generated request ID to ctx.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID returns the ID attached by the HTTP middleware.
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// Write sends a safe public error. Internal causes are only logged.
func Write(w http.ResponseWriter, r *http.Request, status int, cause error) {
	code, message := "internal_error", "Something went wrong. Please try again."
	switch status {
	case http.StatusBadRequest:
		code, message = "invalid_request", "The request is invalid."
	case http.StatusNotFound:
		code, message = "not_found", "The requested resource was not found."
	case http.StatusMethodNotAllowed:
		code, message = "method_not_allowed", "This request method is not allowed."
	}
	id := RequestID(r.Context())
	if status >= 500 {
		slog.ErrorContext(r.Context(), "api error", "request_id", id, "error", cause)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.ErrorResponse{Code: code, Message: message, RequestId: id}); err != nil {
		slog.ErrorContext(r.Context(), "failed to write api error", "request_id", id, "error", err)
	}
}
