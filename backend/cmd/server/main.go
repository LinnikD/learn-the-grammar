package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/LinnikD/learn-the-grammar/backend/internal/api"
	"github.com/LinnikD/learn-the-grammar/backend/internal/apierror"
	"github.com/LinnikD/learn-the-grammar/backend/internal/config"
	"github.com/LinnikD/learn-the-grammar/backend/internal/db"
	"github.com/LinnikD/learn-the-grammar/backend/internal/middleware"
	"github.com/LinnikD/learn-the-grammar/backend/internal/session"
	"github.com/LinnikD/learn-the-grammar/backend/internal/settings"
)

const (
	shutdownTimeout = 10 * time.Second
	sessionTTL      = 30 * 24 * time.Hour
)

// server implements api.StrictServerInterface, the contract generated
// from api/openapi.yaml.
type server struct {
	settings *settings.Service
}

func (server) GetHello(_ context.Context, _ api.GetHelloRequestObject) (api.GetHelloResponseObject, error) {
	return api.GetHello200JSONResponse{Message: "Learn The Grammar!"}, nil
}

func (server) GetMe(ctx context.Context, _ api.GetMeRequestObject) (api.GetMeResponseObject, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		// Unreachable in practice: the Session middleware always sets
		// this before GetMe can run.
		return nil, errors.New("no user id in request context")
	}

	parsed, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("parsing session user id: %w", err)
	}

	return api.GetMe200JSONResponse{UserId: parsed}, nil
}

// GetSettings and PutSettings currently operate on the single stub user
// (db.StubUserID) rather than the request's own session identity — see
// DEBT-1.
func (s server) GetSettings(ctx context.Context, _ api.GetSettingsRequestObject) (api.GetSettingsResponseObject, error) {
	result, err := s.settings.Get(ctx, db.StubUserID)
	if err != nil {
		return nil, fmt.Errorf("getting settings: %w", err)
	}

	return api.GetSettings200JSONResponse(toSettingsResponse(result)), nil
}

func (s server) PutSettings(ctx context.Context, request api.PutSettingsRequestObject) (api.PutSettingsResponseObject, error) {
	if !request.Body.Level.Valid() {
		return api.PutSettingsdefaultJSONResponse{
			Body: api.ErrorResponse{
				Code:      "invalid_request",
				Message:   "The request is invalid.",
				RequestId: apierror.RequestID(ctx),
			},
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	result, err := s.settings.Save(ctx, db.StubUserID, string(request.Body.Level))
	if err != nil {
		return nil, fmt.Errorf("saving settings: %w", err)
	}

	return api.PutSettings200JSONResponse(toSettingsResponse(result)), nil
}

func toSettingsResponse(result settings.Settings) api.SettingsResponse {
	topics := make([]api.Topic, len(result.Topics))
	for i, topic := range result.Topics {
		topics[i] = api.Topic{Id: topic.ID, Name: topic.Name}
	}

	return api.SettingsResponse{
		Level:               api.Level(result.Level),
		OnboardingCompleted: result.OnboardingCompleted,
		Topics:              topics,
	}
}

func newMux(sessionManager *session.Manager, settingsSvc *settings.Service) http.Handler {
	mux := http.NewServeMux()

	badRequest := func(w http.ResponseWriter, r *http.Request, err error) {
		apierror.Write(w, r, http.StatusBadRequest, err)
	}
	strict := api.NewStrictHandlerWithOptions(server{settings: settingsSvc}, nil, api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: badRequest,
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			apierror.Write(w, r, http.StatusInternalServerError, err)
		},
	})
	api.HandlerWithOptions(strict, api.StdHTTPServerOptions{
		ErrorHandlerFunc: badRequest,
		BaseRouter:       mux,
		Middlewares:      []api.MiddlewareFunc{middleware.Session(sessionManager)},
	})

	// /health is infrastructure-only (Kubernetes probes) and deliberately
	// not part of the OpenAPI contract consumed by the frontend, so it
	// stays outside the session middleware too.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Ask ServeMux which route matched. Its built-in 404/405 handlers have
	// no pattern; preserve its Allow header while replacing their text body.
	return middleware.RequestID(middleware.Logging(slog.Default())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, pattern := mux.Handler(r)
		if pattern == "" && (r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/")) {
			result := &routingErrorWriter{header: make(http.Header)}
			handler.ServeHTTP(result, r)
			if allow := result.header.Get("Allow"); allow != "" {
				w.Header().Set("Allow", allow)
			}
			apierror.Write(w, r, result.status, nil)
			return
		}
		mux.ServeHTTP(w, r)
	})))
}

// routingErrorWriter captures only the standard mux's routing error response.
type routingErrorWriter struct {
	header http.Header
	status int
}

func (w *routingErrorWriter) Header() http.Header         { return w.header }
func (w *routingErrorWriter) WriteHeader(status int)      { w.status = status }
func (w *routingErrorWriter) Write(b []byte) (int, error) { return len(b), nil }

// resolveSessionSecret returns cfg's configured secret, or a freshly
// generated random one if none was configured. A random secret is fine
// for local development — it just means sessions don't survive a
// restart — but any long-lived deployment should set LTG_SESSION_SECRET
// explicitly.
func resolveSessionSecret(configured string) (string, error) {
	if configured != "" {
		return configured, nil
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generating random session secret: %w", err)
	}

	slog.Warn("no session secret configured, generated a random one for this run; sessions will not survive a restart")

	return hex.EncodeToString(buf), nil
}

// serve runs an HTTP server on listener until ctx is cancelled, then shuts
// it down gracefully: in-flight requests are given up to shutdownTimeout to
// complete before the server stops.
func serve(ctx context.Context, listener net.Listener, handler http.Handler, shutdownTimeout time.Duration) error {
	srv := &http.Server{Handler: handler}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(listener)
	}()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		slog.Info("shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return srv.Shutdown(shutdownCtx)
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	configPath := flag.String("config", os.Getenv("LTG_CONFIG_FILE"), "path to YAML config file (defaults to LTG_CONFIG_FILE env var)")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if cfg.DatabaseURL == "" {
		slog.Error("LTG_DATABASE_URL is not set")
		os.Exit(1)
	}

	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	settingsSvc := settings.NewService(pool)

	sessionSecret, err := resolveSessionSecret(cfg.SessionSecret)
	if err != nil {
		slog.Error("failed to resolve session secret", "error", err)
		os.Exit(1)
	}
	sessionManager := session.NewManager([]byte(sessionSecret), sessionTTL)

	addr := fmt.Sprintf(":%d", cfg.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("failed to listen", "addr", addr, "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	handler := newMux(sessionManager, settingsSvc)

	slog.Info("server listening", "addr", addr)

	if err := serve(ctx, listener, handler, shutdownTimeout); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
