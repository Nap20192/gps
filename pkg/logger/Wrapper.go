package logger

import (
	"context"
	"log/slog"
)

type HandlerMiddleware struct {
	log slog.Handler
}

func NewHandlerMiddleware(log slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{
		log: log,
	}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.log.Enabled(ctx, rec)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, r slog.Record) error {
	if c, ok := ctx.Value(ctxLoggerKey).(logCtx); ok {

		if c.requestID != "" {
			r.Add(slog.String("request_id", c.requestID))
		}
		if c.action != "" {
			r.Add(slog.String("action", c.action))
		}
		if c.service != "" {
			r.Add(slog.String("service", c.service))
		}
		if c.hostname != "" {
			r.Add(slog.String("hostname", c.hostname))
		}
	}

	if r.Level >= slog.LevelError {
		r.Add(slog.Any("source", r.Source()))
	}
	return h.log.Handle(ctx, r)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.log.WithAttrs(attrs)
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return h.log.WithGroup(name)
}
