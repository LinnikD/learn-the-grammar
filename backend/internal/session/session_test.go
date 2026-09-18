package session

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueVerify_RoundTrip(t *testing.T) {
	m := NewManager([]byte("test-secret"), time.Hour)

	token, err := m.Issue("user-123")
	require.NoError(t, err)

	userID, err := m.Verify(token)
	require.NoError(t, err)
	assert.Equal(t, "user-123", userID)
}

func TestVerify_RejectsExpiredToken(t *testing.T) {
	m := NewManager([]byte("test-secret"), -time.Hour)

	token, err := m.Issue("user-123")
	require.NoError(t, err)

	_, err = m.Verify(token)
	assert.Error(t, err)
}

func TestVerify_RejectsWrongSecret(t *testing.T) {
	issuer := NewManager([]byte("secret-a"), time.Hour)
	verifier := NewManager([]byte("secret-b"), time.Hour)

	token, err := issuer.Issue("user-123")
	require.NoError(t, err)

	_, err = verifier.Verify(token)
	assert.Error(t, err)
}

func TestVerify_RejectsGarbageToken(t *testing.T) {
	m := NewManager([]byte("test-secret"), time.Hour)

	_, err := m.Verify("not-a-jwt")
	assert.Error(t, err)
}

func TestVerify_RejectsUnexpectedSigningMethod(t *testing.T) {
	m := NewManager([]byte("test-secret"), time.Hour)

	// A token signed with "none" (no signature) must never be accepted,
	// regardless of what its claims say.
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
		Subject: "user-123",
	})
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = m.Verify(tokenString)
	assert.Error(t, err)
}
