package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/LinnikD/learn-the-grammar/backend/internal/apierror"
)

// RequestID generates an ID independently of untrusted incoming headers.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(apierror.WithRequestID(r.Context(), id)))
	})
}
