package nexdozlogger

import (
	"context"
	"errors"
	"time"
)

type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Level     string         `json:"level"`
	Service   string         `json:"service"`
	Message   string         `json:"message"`
	TraceID   string         `json:"trace_id,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
	Fields    map[string]any `json:"fields,omitempty"`
}

type Sink interface {
	Write(ctx context.Context, entry Entry) error
	Close() error
}

type Logger struct {
	service string
	min     Level
	sink    Sink
	now     func() time.Time
}

func New(service string, min Level, sink Sink) (*Logger, error) {
	if sink == nil {
		return nil, errors.New("sink is required")
	}
	return &Logger{
		service: service,
		min:     min,
		sink:    sink,
		now:     time.Now,
	}, nil
}

func (l *Logger) Close() error {
	return l.sink.Close()
}

func (l *Logger) Debug(ctx context.Context, message string, fields map[string]any) error {
	return l.log(ctx, LevelDebug, message, fields)
}

func (l *Logger) Info(ctx context.Context, message string, fields map[string]any) error {
	return l.log(ctx, LevelInfo, message, fields)
}

func (l *Logger) Warn(ctx context.Context, message string, fields map[string]any) error {
	return l.log(ctx, LevelWarn, message, fields)
}

func (l *Logger) Error(ctx context.Context, message string, fields map[string]any) error {
	return l.log(ctx, LevelError, message, fields)
}

func (l *Logger) log(ctx context.Context, level Level, message string, fields map[string]any) error {
	if level < l.min {
		return nil
	}

	traceID, requestID := ContextIDs(ctx)
	entry := Entry{
		Timestamp: l.now().UTC(),
		Level:     level.String(),
		Service:   l.service,
		Message:   message,
		TraceID:   traceID,
		RequestID: requestID,
		Fields:    cloneMap(fields),
	}
	return l.sink.Write(ctx, entry)
}

func cloneMap(in map[string]any) map[string]any {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
