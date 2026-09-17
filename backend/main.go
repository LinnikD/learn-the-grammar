package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LinnikD/learn-the-grammar/backend/config"
)

type helloResponse struct {
	Message string `json:"message"`
}

func main() {
	configPath := flag.String("config", os.Getenv("LTG_CONFIG_FILE"), "path to YAML config file (defaults to LTG_CONFIG_FILE env var)")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

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

	addr := fmt.Sprintf(":%d", cfg.Port)

	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
