package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	serverAddr        = os.Getenv("LISTEN_ADDR")
	defaultServerAddr = "localhost:8000"
)

// setupServer creates and configures an HTTP server listening on the address
// specified by the LISTEN_ADDR environment variable or defaulting to port 8000.
// It returns the configured server and the listener for the server to use.
func setupServer(
	ctx context.Context,
	handler http.Handler,
) (*http.Server, net.Listener, error) {
	if serverAddr == "" {
		serverAddr = defaultServerAddr
	}

	listenConfig := net.ListenConfig{
		KeepAlive: 5 * time.Minute,
	}

	listener, err := listenConfig.Listen(ctx, "tcp", serverAddr)
	if err != nil {
		return nil, nil, err
	}

	return &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}, listener, nil
}
