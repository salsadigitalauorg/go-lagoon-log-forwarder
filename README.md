# Go Log Forwarder

A Go package for forwarding logs to UDP endpoints with structured logging using slog.

## Features

- Multi-writer output (stdout and UDP)
- Structured JSON logging
- Configurable via command-line flags
- Built-in log formatting for Logstash/ELK stack

## Usage

```go
import "projects.govcms.gov.au/dev-salsa/go-log-forwarder"

func main() {
    flag.Parse()
    logger.Initialize()

    slog.Info("Application started")
}
```

## Configuration

- `--log.fields.type`: Log type (required) - must match existing k8s namespace
- `--log.host`: UDP host for log forwarding
- `--log.port`: UDP port (default: 5140)
- `--log.addSource`: Add source information to logs (default: true)
- `--log.channel`: Channel name (default: "LagoonLogs")
- `--log.fields.applicationName`: Application name for log identification