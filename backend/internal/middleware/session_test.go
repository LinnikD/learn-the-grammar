package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LinnikD/learn-the-grammar/backend/internal/session"
)

func TestSession_MintsNewSessionWhenNoCookiePresent(t *testing.T) {
	manager := session.NewManager([]byte("test-secret"), time.Hour)

	var gotUserID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := UserIDFromContext(r.Context())
		require.True(t, ok, "expected a user ID in context")
		gotUserID = userID
	})

	handler := Session(manager)(next)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.NotEmpty(t, gotUserID)

	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 1)
	cookie := cookies[0]

	assert.Equal(t, sessionCookieName, cookie.Name)
	assert.True(t, cookie.HttpOnly)
	assert.True(t, cookie.Secure)
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)

	// The cookie must actually verify back to the same user ID that was
	// attached to the request context.
	userID, err := manager.Verify(cookie.Value)
	require.NoError(t, err)
	assert.Equal(t, gotUserID, userID)
}

func TestSession_ReusesValidExistingCookie(t *testing.T) {
	manager := session.NewManager([]byte("test-secret"), time.Hour)

	token, err := manager.Issue("existing-user")
	require.NoError(t, err)

	var gotUserID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID, _ = UserIDFromContext(r.Context())
	})

	handler := Session(manager)(next)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, "existing-user", gotUserID)
	assert.Empty(t, rec.Result().Cookies(), "should not re-issue a cookie for an already-valid session")
}

func TestSession_ReplacesInvalidCookie(t *testing.T) {
	manager := session.NewManager([]byte("test-secret"), time.Hour)

	var gotUserID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID, _ = UserIDFromContext(r.Context())
	})

	handler := Session(manager)(next)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "tampered-or-garbage"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.NotEmpty(t, gotUserID)
	require.Len(t, rec.Result().Cookies(), 1)

	userID, err := manager.Verify(rec.Result().Cookies()[0].Value)
	require.NoError(t, err)
	assert.Equal(t, gotUserID, userID)
}
