package logger

import (
	"fmt"
	"log/slog"
	"os"
)

var Log *slog.Logger

// InitLogger initializes the logger with the specified log level.
// Valid log levels are: "debug", "info", "warn", "error".
func InitLogger(level string) error {
	var logLevel slog.Level

	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	Log = slog.New(handler)
	return nil
}
