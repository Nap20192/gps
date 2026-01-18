package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func InitLogger(level string, pretty bool) (*slog.Logger, error) {
	logLevel := slog.LevelInfo

	switch level {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		fmt.Printf("Unknown log level %s, defaulting to INFO\n", level)
		logLevel = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler

	if pretty {
		handler = NewHandler(WithColor(true), WithLevel(logLevel), WithEncoder(JSON), WithWriter(os.Stdout))
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	logger := slog.New(NewHandlerMiddleware(handler))
	return logger, nil
}
