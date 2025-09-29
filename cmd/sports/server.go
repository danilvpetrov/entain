package main

import (
	"context"
	"net"
	"os"
	"time"

	sportsapi "github.com/danilvpetrov/entain/api/sports"
	"github.com/danilvpetrov/entain/sports"
	"google.golang.org/grpc"
)

var (
	serverAddr        = os.Getenv("LISTEN_ADDR")
	defaultServerAddr = "localhost:9010"
)

// setupServer sets up and returns a gRPC server along with its listener.
func setupServer(
	ctx context.Context,
	s *sports.Service,
) (*grpc.Server, net.Listener, error) {
	server := grpc.NewServer()
	sportsapi.RegisterSportsServer(server, s)

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

	return server, listener, nil
}
