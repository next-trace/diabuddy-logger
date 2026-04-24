package nexdozlogger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const defaultBetterStackEndpoint = "https://in.logs.betterstack.com"

type BetterStackHTTPSink struct {
	client   *http.Client
	endpoint string
	token    string
}

func NewBetterStackHTTPSink(token string, endpoint string) (*BetterStackHTTPSink, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("betterstack source token is required")
	}
	resolvedEndpoint := strings.TrimSpace(endpoint)
	if resolvedEndpoint == "" {
		resolvedEndpoint = defaultBetterStackEndpoint
	}
	return &BetterStackHTTPSink{
		client: &http.Client{Timeout: 4 * time.Second},
		token:  token,
		// Better Stack accepts JSON logs on the endpoint root.
		endpoint: strings.TrimRight(resolvedEndpoint, "/"),
	}, nil
}

func (s *BetterStackHTTPSink) Write(ctx context.Context, entry Entry) error {
	payload, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("betterstack sink returned status %d", res.StatusCode)
	}
	return nil
}

func (s *BetterStackHTTPSink) Close() error {
	return nil
}
