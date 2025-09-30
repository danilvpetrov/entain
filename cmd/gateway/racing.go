package main

import (
	"context"
	"os"

	"github.com/danilvpetrov/entain/api/racing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	racingServiceAddr        = os.Getenv("RACING_SERVICE_ADDR")
	defaultRacingServiceAddr = "localhost:9000"
)

// setupRacingService sets up the gRPC gateway for the Racing service, allowing
// HTTP requests to be proxied to the gRPC server.
func setupRacingService(ctx context.Context, mux *runtime.ServeMux) error {
	if racingServiceAddr == "" {
		racingServiceAddr = defaultRacingServiceAddr
	}

	otelClientHdr := otelgrpc.NewClientHandler()

	return racing.RegisterRacingHandlerFromEndpoint(
		ctx,
		mux,
		racingServiceAddr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
			grpc.WithStatsHandler(otelClientHdr),
		},
	)
}
