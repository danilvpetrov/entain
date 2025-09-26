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

	mux, err := setupAPI(ctx)
	if err != nil {
		return fmt.Errorf("error setting up API: %w", err)
	}

	svr, listener, err := setupServer(mux)
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
		// Do not use context from main, as it is already cancelled.
		if err := svr.Shutdown(context.Background()); err != nil {
			slog.Error("error shutting down server", "error", err)
		}
	}()

	slog.Info(
		"gateway server listening",
		slog.String("addr", listener.Addr().String()),
	)

	return svr.Serve(listener)
}
