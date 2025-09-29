package main

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// setupAPI sets up the HTTP API gateway, routing requests to the appropriate
// gRPC services.
func setupAPI(ctx context.Context) (*runtime.ServeMux, error) {
	m := runtime.NewServeMux()

	if err := setupRacingService(ctx, m); err != nil {
		return nil, fmt.Errorf("error setting up racing service: %w", err)
	}

	if err := setupSportsService(ctx, m); err != nil {
		return nil, fmt.Errorf("error setting up sports service: %w", err)
	}

	return m, nil
}
