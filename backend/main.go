package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type helloResponse struct {
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(helloResponse{
			Message: "Hello World",
		}); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	})

	log.Println("server listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}