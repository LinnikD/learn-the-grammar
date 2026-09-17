package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServe_GracefulShutdownWaitsForInFlightRequest(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})

	mux := http.NewServeMux()
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		close(requestStarted)
		<-releaseRequest
		w.WriteHeader(http.StatusOK)
	})

	ctx, cancel := context.WithCancel(context.Background())

	serveDone := make(chan error, 1)
	go func() {
		serveDone <- serve(ctx, listener, mux, 2*time.Second)
	}()

	requestDone := make(chan error, 1)
	go func() {
		resp, err := http.Get("http://" + listener.Addr().String() + "/slow")
		if err != nil {
			requestDone <- err
			return
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			requestDone <- fmt.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
			return
		}
		requestDone <- nil
	}()

	select {
	case <-requestStarted:
	case <-time.After(time.Second):
		t.Fatal("request never reached the handler")
	}

	// Trigger shutdown while the request is still in flight.
	cancel()

	// Give shutdown a moment to start, then let the handler finish. If
	// shutdown were not graceful, the connection would already be gone by
	// the time we release the handler and the client request would fail.
	time.Sleep(50 * time.Millisecond)
	close(releaseRequest)

	select {
	case err := <-requestDone:
		if err != nil {
			t.Fatalf("in-flight request failed during shutdown: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("in-flight request did not complete")
	}

	select {
	case err := <-serveDone:
		if err != nil {
			t.Fatalf("serve returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("serve did not return after shutdown")
	}
}

func TestServe_ReturnsPromptlyWithNoInFlightRequests(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)
	go func() {
		done <- serve(ctx, listener, http.NewServeMux(), 2*time.Second)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("serve returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("serve did not return promptly for an already-cancelled context")
	}
}
