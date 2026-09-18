package main

import (
	"context"
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

	"github.com/LinnikD/learn-the-grammar/backend/internal/api"
	"github.com/LinnikD/learn-the-grammar/backend/internal/config"
	"github.com/LinnikD/learn-the-grammar/backend/internal/middleware"
)

const shutdownTimeout = 10 * time.Second

// server implements api.StrictServerInterface, the contract generated
// from api/openapi.yaml.
type server struct{}

func (server) GetHello(_ context.Context, _ api.GetHelloRequestObject) (api.GetHelloResponseObject, error) {
	return api.GetHello200JSONResponse{Message: "Learn The Grammar!"}, nil
}

func newMux() http.Handler {
	mux := http.NewServeMux()

	api.HandlerFromMux(api.NewStrictHandler(server{}, nil), mux)

	// /health is infrastructure-only (Kubernetes probes) and deliberately
	// not part of the OpenAPI contract consumed by the frontend.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return mux
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

	addr := fmt.Sprintf(":%d", cfg.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("failed to listen", "addr", addr, "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	handler := middleware.Logging(logger)(newMux())

	slog.Info("server listening", "addr", addr)

	if err := serve(ctx, listener, handler, shutdownTimeout); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
