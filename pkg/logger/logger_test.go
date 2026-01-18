package logger

import (
	"log/slog"
	"testing"
)

func TestLogger(t *testing.T) {
	logger, err := InitLogger("debug", true)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	slog.SetDefault(logger)

	slog.Debug("This is a debug message", slog.String("key1", "value1"))
	slog.Info("This is an info message", slog.Int("key2", 42))
	slog.Warn("This is a warning message", slog.Float64("key3", 3.14))
	slog.Error("This is an error message", slog.Bool("key4", true))
}
