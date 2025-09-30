package main

import (
	"context"
	"os"

	"github.com/danilvpetrov/entain/api/sports"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	sportsServiceAddr        = os.Getenv("SPORTS_SERVICE_ADDR")
	defaultSportsServiceAddr = "localhost:9010"
)

// setupSportsService sets up the gRPC gateway for the Sports service, allowing
// HTTP requests to be proxied to the gRPC server.
func setupSportsService(ctx context.Context, mux *runtime.ServeMux) error {
	if sportsServiceAddr == "" {
		sportsServiceAddr = defaultSportsServiceAddr
	}

	otelClientHdr := otelgrpc.NewClientHandler()

	return sports.RegisterSportsHandlerFromEndpoint(
		ctx,
		mux,
		sportsServiceAddr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
			grpc.WithStatsHandler(otelClientHdr),
		},
	)
}
