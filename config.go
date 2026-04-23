package diabuddylogger

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Service            string
	Level              Level
	Sink               string
	BetterStackToken   string
	BetterStackAddress string
}

func ConfigFromEnv(service string) (Config, error) {
	level, err := ParseLevel(os.Getenv("LOGGER_LEVEL"))
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		Service:            strings.TrimSpace(service),
		Level:              level,
		Sink:               strings.TrimSpace(strings.ToLower(os.Getenv("LOGGER_SINK"))),
		BetterStackToken:   strings.TrimSpace(os.Getenv("LOGGER_BETTERSTACK_SOURCE_TOKEN")),
		BetterStackAddress: strings.TrimSpace(os.Getenv("LOGGER_BETTERSTACK_ENDPOINT")),
	}
	if cfg.Sink == "" {
		cfg.Sink = "stdout"
	}
	return cfg, nil
}

func NewFromConfig(cfg Config) (*Logger, error) {
	var sink Sink
	switch cfg.Sink {
	case "stdout":
		sink = NewStdoutSink()
	case "noop":
		sink = NewNoopSink()
	case "betterstack":
		httpSink, err := NewBetterStackHTTPSink(cfg.BetterStackToken, cfg.BetterStackAddress)
		if err != nil {
			return nil, err
		}
		sink = httpSink
	default:
		return nil, fmt.Errorf("unsupported LOGGER_SINK: %s", cfg.Sink)
	}
	return New(cfg.Service, cfg.Level, sink)
}

func NewFromEnv(service string) (*Logger, error) {
	cfg, err := ConfigFromEnv(service)
	if err != nil {
		return nil, err
	}
	return NewFromConfig(cfg)
}
