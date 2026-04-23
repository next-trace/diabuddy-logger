# diabuddy-logger

Minimal, pluggable logger for DiaBuddy services.

## Sinks

- `stdout` (default): structured JSON to stdout.
- `betterstack`: HTTP ingestion to Better Stack (Logtail endpoint).
- `noop`: disables writes.

## Environment

- `LOGGER_SINK=stdout|betterstack|noop`
- `LOGGER_LEVEL=debug|info|warn|error`
- `LOGGER_BETTERSTACK_SOURCE_TOKEN=...` (required for `betterstack`)
- `LOGGER_BETTERSTACK_ENDPOINT=https://in.logs.betterstack.com` (optional)

## Usage

```go
logger, err := diabuddylogger.NewFromEnv("diabuddy-user-api")
if err != nil {
  panic(err)
}
defer logger.Close()

ctx := diabuddylogger.WithRequestID(context.Background(), "req-123")
_ = logger.Info(ctx, "request completed", map[string]any{
  "path": "/healthz",
})
```

## DigitalOcean note

DigitalOcean does not currently provide a full free managed log platform comparable to a dedicated SaaS logs product.
For MVP/dev cost control, use one of:

- App/Container stdout logs (short retention, basic troubleshooting)
- Better Stack free tier via `betterstack` sink
- Self-hosted Loki on a small Droplet (low-cost, ops overhead)
