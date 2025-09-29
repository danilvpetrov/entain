package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

var dbg = os.Getenv("DEBUG")

// setupLogger configures the global logger based on the DEBUG environment
// variable. If DEBUG is set to "true", the logger will output debug-level logs
// in a human-readable text format. Otherwise, it will log in JSON format with
// the default log level.
func setupLogger() error {
	var (
		isDbg bool
		err   error
	)

	if dbg != "" {
		isDbg, err = strconv.ParseBool(dbg)
		if err != nil {
			return fmt.Errorf("error parsing DEBUG envvar: %w", err)
		}
	}

	var logger *slog.Logger
	if isDbg {
		logger = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)

	return nil
}
