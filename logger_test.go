package nexdozlogger

import (
	"context"
	"testing"
	"time"
)

type captureSink struct {
	entries []Entry
}

func (s *captureSink) Write(_ context.Context, entry Entry) error {
	s.entries = append(s.entries, entry)
	return nil
}

func (s *captureSink) Close() error {
	return nil
}

func TestParseLevel(t *testing.T) {
	level, err := ParseLevel("debug")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if level != LevelDebug {
		t.Fatalf("expected debug level, got %v", level)
	}

	if _, err := ParseLevel("broken"); err == nil {
		t.Fatalf("expected error for invalid level")
	}
}

func TestLoggerMinLevel(t *testing.T) {
	sink := &captureSink{}
	logger, err := New("user-api", LevelInfo, sink)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	logger.now = func() time.Time { return time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC) }

	if err := logger.Debug(context.Background(), "hidden", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sink.entries) != 0 {
		t.Fatalf("expected no entries for debug below info")
	}

	if err := logger.Info(context.Background(), "visible", map[string]any{"k": "v"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sink.entries) != 1 {
		t.Fatalf("expected one entry, got %d", len(sink.entries))
	}
	if sink.entries[0].Message != "visible" {
		t.Fatalf("unexpected message: %s", sink.entries[0].Message)
	}
	if sink.entries[0].Fields["k"] != "v" {
		t.Fatalf("expected field k=v")
	}
}

func TestNewFromConfig(t *testing.T) {
	logger, err := NewFromConfig(Config{Service: "user-api", Level: LevelInfo, Sink: "noop"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if logger == nil {
		t.Fatalf("expected logger")
	}

	if _, err := NewFromConfig(Config{Service: "user-api", Level: LevelInfo, Sink: "missing"}); err == nil {
		t.Fatalf("expected error for unsupported sink")
	}
}
