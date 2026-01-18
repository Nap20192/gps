package logger

import (
	"context"
)

type contextKey struct{}

var ctxLoggerKey = contextKey{}

type logCtx struct {
	requestID string
	action    string
	service   string
	hostname  string
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	c := getOrCreateLogCtx(ctx)
	c.requestID = requestID
	return context.WithValue(ctx, ctxLoggerKey, c)
}

func WithAction(ctx context.Context, action string) context.Context {
	c := getOrCreateLogCtx(ctx)
	c.action = action
	return context.WithValue(ctx, ctxLoggerKey, c)
}

func WithService(ctx context.Context, service string) context.Context {
	c := getOrCreateLogCtx(ctx)
	c.service = service
	return context.WithValue(ctx, ctxLoggerKey, c)
}

func WithHostname(ctx context.Context, hostname string) context.Context {
	c := getOrCreateLogCtx(ctx)
	c.hostname = hostname
	return context.WithValue(ctx, ctxLoggerKey, c)
}

func getOrCreateLogCtx(ctx context.Context) logCtx {
	if c, ok := ctx.Value(ctxLoggerKey).(logCtx); ok {
		return c
	}
	return logCtx{}
}

func WithSessionId(ctx context.Context, sessionId string) context.Context {
	return WithRequestID(ctx, sessionId)
}

func WithName(ctx context.Context, name string) context.Context {
	return WithAction(ctx, name)
}
