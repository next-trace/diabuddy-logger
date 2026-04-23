package diabuddylogger

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"
)

type StdoutSink struct {
	out io.Writer
	mu  sync.Mutex
}

func NewStdoutSink() *StdoutSink {
	return &StdoutSink{out: os.Stdout}
}

func (s *StdoutSink) Write(_ context.Context, entry Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = s.out.Write(append(b, '\n'))
	return err
}

func (s *StdoutSink) Close() error {
	return nil
}
