package diabuddylogger

import "context"

type NoopSink struct{}

func NewNoopSink() *NoopSink {
	return &NoopSink{}
}

func (s *NoopSink) Write(_ context.Context, _ Entry) error {
	return nil
}

func (s *NoopSink) Close() error {
	return nil
}
