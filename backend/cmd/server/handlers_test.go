package main

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LinnikD/learn-the-grammar/backend/internal/session"
)

func TestGetMe_ReturnsAndPersistsSessionUserID(t *testing.T) {
	manager := session.NewManager([]byte("test-secret"), time.Hour)
	handler := newMux(manager)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	client := srv.Client()
	client.Jar = jar

	// First request: no session yet, one should be minted.
	firstUserID := requestMe(t, client, srv.URL)
	assert.NotEmpty(t, firstUserID)

	// Second request, same client (same cookie jar): must be the same user.
	secondUserID := requestMe(t, client, srv.URL)
	assert.Equal(t, firstUserID, secondUserID)

	// A different client (no cookie jar / no cookie sent) gets a
	// different, freshly minted user.
	otherUserID := requestMe(t, http.DefaultClient, srv.URL)
	assert.NotEqual(t, firstUserID, otherUserID)
}

func requestMe(t *testing.T, client *http.Client, baseURL string) string {
	t.Helper()

	resp, err := client.Get(baseURL + "/api/me")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var body struct {
		UserID string `json:"user_id"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

	return body.UserID
}

func TestResolveSessionSecret_UsesConfiguredValue(t *testing.T) {
	secret, err := resolveSessionSecret("configured-secret")
	require.NoError(t, err)
	assert.Equal(t, "configured-secret", secret)
}

func TestResolveSessionSecret_GeneratesRandomWhenEmpty(t *testing.T) {
	first, err := resolveSessionSecret("")
	require.NoError(t, err)
	assert.NotEmpty(t, first)

	second, err := resolveSessionSecret("")
	require.NoError(t, err)

	assert.NotEqual(t, first, second, "expected two independently generated secrets to differ")
}

func TestAPIRoutingErrors(t *testing.T) {
	handler := newMux(session.NewManager([]byte("test-secret"), time.Hour))
	for _, tc := range []struct {
		method, path string
		status       int
		code, allow  string
	}{
		{"GET", "/api/missing", 404, "not_found", ""},
		{"GET", "/api", 404, "not_found", ""},
		{"POST", "/api/me", 405, "method_not_allowed", "GET, HEAD"},
		{"POST", "/api/hello", 405, "method_not_allowed", "GET, HEAD"},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("X-Request-ID", "untrusted-client-id")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, tc.status, rec.Code)
			assert.Equal(t, tc.allow, rec.Header().Get("Allow"))
			var body struct {
				Code      string
				RequestID string `json:"request_id"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tc.code, body.Code)
			assert.NotEmpty(t, body.RequestID)
			assert.NotEqual(t, "untrusted-client-id", body.RequestID)
			assert.Equal(t, rec.Header().Get("X-Request-ID"), body.RequestID)
		})
	}
}

func TestHandlerErrorUsesPublicContract(t *testing.T) {
	manager := session.NewManager([]byte("test-secret"), time.Hour)
	// A signed but malformed subject reaches the real GetMe error handler.
	token, err := manager.Issue("private-invalid-user-id")
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(&http.Cookie{Name: "ltg_session", Value: token})
	rec := httptest.NewRecorder()
	newMux(manager).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	var body struct {
		Code, Message string
		RequestID     string `json:"request_id"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "internal_error", body.Code)
	assert.Equal(t, "Something went wrong. Please try again.", body.Message)
	assert.NotContains(t, rec.Body.String(), "private-invalid-user-id")
	assert.Equal(t, rec.Header().Get("X-Request-ID"), body.RequestID)
}
