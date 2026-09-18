package apierror

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LinnikD/learn-the-grammar/backend/internal/api"
)

func TestWrite(t *testing.T) {
	var logs bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(original) })
	for _, tc := range []struct {
		status int
		code   string
	}{
		{400, "invalid_request"}, {404, "not_found"},
		{405, "method_not_allowed"}, {500, "internal_error"},
	} {
		t.Run(tc.code, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
			req = req.WithContext(WithRequestID(req.Context(), "test-id"))
			rec := httptest.NewRecorder()
			Write(rec, req, tc.status, errors.New("private database detail"))
			assert.Equal(t, tc.status, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			var body api.ErrorResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tc.code, body.Code)
			assert.Equal(t, "test-id", body.RequestId)
			assert.NotEmpty(t, body.Message)
			assert.NotContains(t, rec.Body.String(), "private database detail")
		})
	}
	var entry map[string]any
	require.NoError(t, json.Unmarshal(logs.Bytes(), &entry))
	assert.Equal(t, "test-id", entry["request_id"])
	assert.Equal(t, "private database detail", entry["error"])
}
