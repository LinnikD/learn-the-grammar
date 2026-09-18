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
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/LinnikD/learn-the-grammar/backend/internal/api"
	"github.com/LinnikD/learn-the-grammar/backend/internal/config"
	"github.com/LinnikD/learn-the-grammar/backend/internal/middleware"
	"github.com/LinnikD/learn-the-grammar/backend/internal/session"
)

const (
	shutdownTimeout = 10 * time.Second
	sessionTTL      = 30 * 24 * time.Hour
)

// server implements api.StrictServerInterface, the contract generated
// from api/openapi.yaml.
type server struct{}

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

func newMux(sessionManager *session.Manager) http.Handler {
	mux := http.NewServeMux()

	api.HandlerWithOptions(api.NewStrictHandler(server{}, nil), api.StdHTTPServerOptions{
		BaseRouter:  mux,
		Middlewares: []api.MiddlewareFunc{middleware.Session(sessionManager)},
	})

	// /health is infrastructure-only (Kubernetes probes) and deliberately
	// not part of the OpenAPI contract consumed by the frontend, so it
	// stays outside the session middleware too.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return mux
}

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

	handler := middleware.Logging(logger)(newMux(sessionManager))

	slog.Info("server listening", "addr", addr)

	if err := serve(ctx, listener, handler, shutdownTimeout); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
