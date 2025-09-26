package main

import (
	"net"
	"net/http"
	"os"
	"time"
)

var (
	serverAddr        = os.Getenv("LISTEN_ADDR")
	defaultServerAddr = ":8000"
)

// setupServer creates and configures an HTTP server listening on the address
// specified by the LISTEN_ADDR environment variable or defaulting to port 8000.
// It returns the configured server and the listener for the server to use.
func setupServer(handler http.Handler) (*http.Server, net.Listener, error) {
	if serverAddr == "" {
		serverAddr = defaultServerAddr
	}

	l, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return nil, nil, err
	}

	return &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}, l, nil
}
