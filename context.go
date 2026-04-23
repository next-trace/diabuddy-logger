package diabuddylogger

import "context"

type contextKey string

const (
	traceIDKey   contextKey = "diabuddy.trace_id"
	requestIDKey contextKey = "diabuddy.request_id"
)

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func ContextIDs(ctx context.Context) (traceID string, requestID string) {
	if ctx == nil {
		return "", ""
	}
	if value, ok := ctx.Value(traceIDKey).(string); ok {
		traceID = value
	}
	if value, ok := ctx.Value(requestIDKey).(string); ok {
		requestID = value
	}
	return traceID, requestID
}
