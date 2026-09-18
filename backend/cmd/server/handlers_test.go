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
