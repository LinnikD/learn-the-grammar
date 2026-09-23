package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LinnikD/learn-the-grammar/backend/internal/db/dbtest"
	"github.com/LinnikD/learn-the-grammar/backend/internal/session"
	"github.com/LinnikD/learn-the-grammar/backend/internal/settings"
)

func newSettingsTestHandler(t *testing.T) http.Handler {
	t.Helper()

	pool := dbtest.StartPostgres(t)
	manager := session.NewManager([]byte("test-secret"), time.Hour)

	return newMux(manager, settings.NewService(pool))
}

func TestSettings_GetThenPut(t *testing.T) {
	handler := newSettingsTestHandler(t)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	client := srv.Client()

	var got struct {
		Level               string `json:"level"`
		OnboardingCompleted bool   `json:"onboarding_completed"`
		Topics              []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"topics"`
	}

	resp, err := client.Get(srv.URL + "/api/settings")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.NoError(t, resp.Body.Close())

	assert.Equal(t, "A1", got.Level)
	assert.False(t, got.OnboardingCompleted)

	body, err := json.Marshal(map[string]string{"level": "B1"})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, srv.URL+"/api/settings", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&got))
	require.NoError(t, resp2.Body.Close())

	assert.Equal(t, "B1", got.Level)
	assert.True(t, got.OnboardingCompleted)
}

func TestSettings_PutRejectsInvalidLevel(t *testing.T) {
	handler := newSettingsTestHandler(t)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	body, err := json.Marshal(map[string]string{"level": "not-a-level"})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, srv.URL+"/api/settings", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := srv.Client().Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errBody struct {
		Code string `json:"code"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&errBody))
	assert.Equal(t, "invalid_request", errBody.Code)
}
