package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/LinnikD/learn-the-grammar/backend/internal/session"
)

const sessionCookieName = "ltg_session"

type contextKey int

const userIDContextKey contextKey = iota

// UserIDFromContext returns the user ID attached to ctx by Session, and
// whether one was present.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}

// Session ensures every request carries a valid session. If the
// request's session cookie is missing or fails verification (invalid
// signature, expired), a new one is minted for a freshly generated user
// ID and set on the response — there is no login yet, so any visitor
// is treated as a new user until one exists. Either way, the resulting
// user ID is attached to the request context for downstream handlers.
func Session(manager *session.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := verifySessionCookie(manager, r)
			if !ok {
				userID = uuid.NewString()

				if err := issueSessionCookie(manager, w, userID); err != nil {
					http.Error(w, "failed to create session", http.StatusInternalServerError)
					return
				}
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func verifySessionCookie(manager *session.Manager, r *http.Request) (string, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return "", false
	}

	userID, err := manager.Verify(cookie.Value)
	if err != nil {
		return "", false
	}

	return userID, true
}

func issueSessionCookie(manager *session.Manager, w http.ResponseWriter, userID string) error {
	token, err := manager.Issue(userID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(manager.TTL().Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}
