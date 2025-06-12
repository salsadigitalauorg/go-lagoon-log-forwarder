# Go Lagoon Log Forwarder

[![CI](https://github.com/salsadigitalauorg/go-lagoon-log-forwarder/actions/workflows/pr.yml/badge.svg)](https://github.com/salsadigitalauorg/go-lagoon-log-forwarder/actions/workflows/pr.yml)
[![Release](https://github.com/salsadigitalauorg/go-lagoon-log-forwarder/actions/workflows/release.yml/badge.svg)](https://github.com/salsadigitalauorg/go-lagoon-log-forwarder/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/salsadigitalauorg/go-lagoon-log-forwarder/branch/main/graph/badge.svg)](https://codecov.io/gh/salsadigitalauorg/go-lagoon-log-forwarder)
[![Go Report Card](https://goreportcard.com/badge/github.com/salsadigitalauorg/go-lagoon-log-forwarder)](https://goreportcard.com/report/github.com/salsadigitalauorg/go-lagoon-log-forwarder)
[![Go Reference](https://pkg.go.dev/badge/github.com/salsadigitalauorg/go-lagoon-log-forwarder.svg)](https://pkg.go.dev/github.com/salsadigitalauorg/go-lagoon-log-forwarder)

A high-performance Go library for forwarding structured logs to UDP endpoints, designed for Lagoon/Kubernetes environments with built-in support for ELK stack integration.

## ✨ Features

- **🚀 High Performance**: Multi-writer output (stdout + UDP) with efficient connection handling
- **📋 Structured Logging**: JSON logging using Go's native `slog` package
- **⚙️ Flexible Configuration**: Programmatic configuration (no global flags)
- **🔧 ELK Stack Ready**: Built-in log formatting for Logstash/Elasticsearch
- **🛡️ Production Ready**: Comprehensive error handling and graceful failures
- **🧪 Well Tested**: 100% test coverage with benchmarks
- **📦 Zero Dependencies**: Uses only Go standard library

## 📦 Installation

```bash
go get github.com/salsadigitalauorg/go-lagoon-log-forwarder@latest
```

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
    "log/slog"
    
    "github.com/salsadigitalauorg/go-lagoon-log-forwarder"
)

func main() {
    // Create configuration with defaults
    cfg := logger.NewConfig()
    cfg.LogType = "my-app"        // Required: must match k8s namespace
    cfg.LogHost = "logstash.example.com"
    
    // Initialize logger
    if err := logger.Initialize(cfg); err != nil {
        panic(err)
    }
    
    // Use structured logging
    slog.Info("Application started", 
        "version", "1.0.0",
        "environment", "production",
    )
    
    slog.Error("Database connection failed", 
        "error", err,
        "database", "postgres",
        "retry_count", 3,
    )
}
```

### Advanced Configuration

```go
cfg := logger.NewConfig()
cfg.LogType = "my-microservice"
cfg.LogHost = "logs.k8s.cluster"
cfg.LogPort = 5140
cfg.ApplicationName = "user-service"
cfg.LogChannel = "ProductionLogs"
cfg.AddSource = true  // Include source file/line info

if err := logger.Initialize(cfg); err != nil {
    log.Fatalf("Failed to initialize logger: %v", err)
}
```

## ⚙️ Configuration Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `LogType` | `string` | **required** | Log type (must match k8s namespace) |
| `LogHost` | `string` | `""` | UDP host for log forwarding |
| `LogPort` | `int` | `5140` | UDP port number |
| `ApplicationName` | `string` | `""` | Application identifier |
| `LogChannel` | `string` | `"LagoonLogs"` | Channel name for log routing |
| `AddSource` | `bool` | `true` | Include source file/line information |
| `MessageVersion` | `int` | `1` | Log message format version |

## 📝 Log Format

The logger produces structured JSON logs compatible with ELK stack:

```json
{
  "time": "2024-01-15T10:30:00Z",
  "level": "INFO", 
  "message": "User authenticated",
  "@version": 3,
  "@timestamp": "2024-01-15T10:30:00Z",
  "application": "user-service",
  "channel": "LagoonLogs",
  "host": "pod-abc123",
  "type": "production",
  "user_id": 12345,
  "method": "POST"
}
```

### Automatic Field Mapping

- `msg` → `message`
- `time` → `@timestamp` 
- `timestampOverride` → `@timestamp`

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────┐    ┌─────────────────┐
│   Application   │    │   Logger     │    │   Multi-Writer  │
│                 │───▶│              │───▶│                 │
│  slog.Info()    │    │  Transform   │    │  ┌─────────────┐│
│  slog.Error()   │    │  Fields      │    │  │   Stdout    ││
└─────────────────┘    └──────────────┘    │  └─────────────┘│
                                           │  ┌─────────────┐│
                                           │  │ UDP Socket  ││
                                           │  │ (Logstash)  ││
                                           │  └─────────────┘│
                                           └─────────────────┘
```

## 🧪 Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem ./...
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or later
- Git

### Local Setup

```bash
git clone https://github.com/salsadigitalauorg/go-lagoon-log-forwarder.git
cd go-lagoon-log-forwarder

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run linter (requires golangci-lint)
golangci-lint run
```

### Code Quality

This project maintains high code quality standards:

- **100% Test Coverage**: All code paths tested
- **Comprehensive Linting**: 30+ linters via golangci-lint
- **Security Scanning**: Gosec security analysis
- **Performance Testing**: Benchmark tests included
- **Multi-Version Support**: Tested on Go 1.21-1.24

## 🚀 CI/CD Pipeline

This project uses automated CI/CD with:

- **Pull Request Testing**: Comprehensive checks on all PRs
- **Automatic Versioning**: Semantic versioning based on conventional commits
- **Automated Releases**: GitHub releases with auto-generated changelogs
- **Cross-Platform Testing**: Linux, macOS, Windows compatibility

### Commit Message Format

Use [Conventional Commits](https://conventionalcommits.org/) for automatic versioning:

```bash
feat: add custom log formatter support       # Minor version bump
fix: resolve UDP connection memory leak      # Patch version bump  
feat!: change Initialize function signature  # Major version bump
```

## 📊 Performance

Benchmark results on typical hardware:

```
BenchmarkNewConfig-10      1000000000    0.23 ns/op    0 B/op    0 allocs/op
BenchmarkConfig-10            2441269  475.40 ns/op   48 B/op    1 allocs/op
BenchmarkDefaultAttrs-10     50000000   28.50 ns/op   24 B/op    1 allocs/op
BenchmarkReplaceAttr-10     100000000   12.30 ns/op    0 B/op    0 allocs/op
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:

- Development workflow
- Commit message conventions
- Code quality standards
- Testing requirements

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🏷️ Versioning

This project uses [Semantic Versioning](https://semver.org/):

- **Major**: Breaking changes
- **Minor**: New features (backward compatible)
- **Patch**: Bug fixes and improvements

Releases are automated based on conventional commit messages.

## 📞 Support

- **🐛 Bug Reports**: [GitHub Issues](https://github.com/salsadigitalauorg/go-lagoon-log-forwarder/issues)
- **💡 Feature Requests**: [GitHub Issues](https://github.com/salsadigitalauorg/go-lagoon-log-forwarder/issues)
- **❓ Questions**: [GitHub Discussions](https://github.com/salsadigitalauorg/go-lagoon-log-forwarder/discussions)

---

**Built with ❤️ by the Salsa Digital team for Lagoon & Kubernetes deployments**