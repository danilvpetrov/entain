package main

import (
	"net"
	"os"

	racingapi "github.com/danilvpetrov/entain/api/racing"
	"github.com/danilvpetrov/entain/racing"
	"google.golang.org/grpc"
)

var (
	serverAddr        = os.Getenv("LISTEN_ADDR")
	defaultServerAddr = ":9000"
)

// setupServer sets up and returns a gRPC server along with its listener.
func setupServer(s *racing.Service) (*grpc.Server, net.Listener, error) {
	server := grpc.NewServer()
	racingapi.RegisterRacingServer(server, s)

	if serverAddr == "" {
		serverAddr = defaultServerAddr
	}

	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return nil, nil, err
	}

	return server, listener, nil
}
