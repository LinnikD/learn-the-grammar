package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LinnikD/learn-the-grammar/backend/internal/config"
)

const shutdownTimeout = 10 * time.Second

type helloResponse struct {
	Message string `json:"message"`
}

func newMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(helloResponse{
			Message: "Learn The Grammar!",
		}); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	})

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
		log.Println("shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return srv.Shutdown(shutdownCtx)
	}
}

func main() {
	configPath := flag.String("config", os.Getenv("LTG_CONFIG_FILE"), "path to YAML config file (defaults to LTG_CONFIG_FILE env var)")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("server listening on %s", addr)

	if err := serve(ctx, listener, newMux(), shutdownTimeout); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
