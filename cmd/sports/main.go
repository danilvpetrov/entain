package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt, os.Kill,
	)
	defer cancel()

	if err := setupLogger(); err != nil {
		return fmt.Errorf("error setting up logger: %w", err)
	}

	shutdown, err := openTelemetry(ctx)
	if err != nil {
		return fmt.Errorf("error setting up OpenTelemetry: %w", err)
	}
	defer func() {
		if err := shutdown(); err != nil {
			slog.Error(
				"error shutting down OpenTelemetry",
				slog.Any("error", err),
			)
		}
	}()

	db, err := setupDB(ctx)
	if err != nil {
		return fmt.Errorf("error setting up database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing database: %v\n", err)
		}
	}()

	service := setupService(db)
	svr, listener, err := setupServer(ctx, service)
	if err != nil {
		return fmt.Errorf("error setting up server: %w", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			slog.Error("error closing listener", slog.Any("error", err))
		}
	}()

	go func() {
		<-ctx.Done()
		slog.Info("shutting down server")
		svr.GracefulStop()
	}()

	slog.Info(
		"sports server listening",
		slog.String("addr", listener.Addr().String()),
	)

	return svr.Serve(listener)
}
