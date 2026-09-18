// Package session issues and verifies signed session tokens.
//
// A token identifies a user by ID (the JWT subject) and is
// self-contained: no server-side session store exists. There is no
// login yet, so today every user ID is auto-generated on first visit
// (see internal/middleware.Session) rather than backed by a registered
// account — but the token format is meant to carry over once real
// authentication exists, at which point login would issue the same
// kind of token for a real, persisted user ID instead.
package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var errNoSubject = errors.New("token missing subject claim")

// Manager issues and verifies session tokens signed with an HMAC secret.
type Manager struct {
	secret []byte
	ttl    time.Duration
}

// NewManager returns a Manager that signs tokens with secret and issues
// them with the given time-to-live.
func NewManager(secret []byte, ttl time.Duration) *Manager {
	return &Manager{secret: secret, ttl: ttl}
}

// TTL returns the configured token lifetime.
func (m *Manager) TTL() time.Duration {
	return m.ttl
}

// Issue returns a signed token identifying userID, valid for the
// manager's configured TTL.
func (m *Manager) Issue(userID string) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(m.secret)
}

// Verify validates tokenString and returns the user ID it identifies.
func (m *Manager) Verify(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method)
		}

		return m.secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("verifying token: %w", err)
	}

	if claims.Subject == "" {
		return "", errNoSubject
	}

	return claims.Subject, nil
}
