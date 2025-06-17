# Contributing to Go Log Forwarder

Thank you for contributing to Go Log Forwarder! This document explains our automated CI/CD process and development workflow.

## Development Workflow

### 1. Pull Request Process

When you create a Pull Request against the `main` branch:

- **Automated Testing**: The PR workflow runs comprehensive tests across multiple Go versions (1.21-1.24)
- **Code Quality**: Linting, security scanning, and code formatting checks are performed
- **Coverage**: Test coverage is measured and reported to Codecov

### 2. Merge to Main

When a PR is merged to `main`:

- **Automatic Versioning**: A new version is automatically created based on commit messages
- **Release Creation**: A GitHub release is published with changelog
- **Module Publishing**: The Go module is available at the new version tag

## Commit Message Format

We use **Conventional Commits** for automatic versioning. Your commit messages should follow this format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Commit Types

| Type | Description | Version Bump |
|------|-------------|--------------|
| `feat` | New features | Minor |
| `fix` | Bug fixes | Patch |
| `perf` | Performance improvements | Patch |
| `refactor` | Code refactoring | Patch |
| `docs` | Documentation changes | Patch |
| `test` | Test additions/changes | None |
| `chore` | Maintenance tasks | None |
| `ci` | CI/CD changes | None |
| `build` | Build system changes | None |

### Examples

```bash
# New feature (bumps minor version: 1.0.0 → 1.1.0)
feat: add support for custom log formatters

# Bug fix (bumps patch version: 1.1.0 → 1.1.1)  
fix: resolve memory leak in UDP connection pooling

# Breaking change (bumps major version: 1.1.1 → 2.0.0)
feat!: change Initialize function signature

# With scope
feat(config): add environment variable support
fix(logger): handle nil pointer in defaultAttrs
```

### Breaking Changes

For breaking changes, add `!` after the type or include `BREAKING CHANGE:` in the footer:

```bash
feat!: change config struct field names

# OR

feat: redesign configuration system

BREAKING CHANGE: Config field names have changed from snake_case to camelCase
```

## Local Development

### Prerequisites

- Go 1.21 or later
- Git
- Make (for using the Makefile commands)

### Setup

```bash
# Clone the repository
git clone https://github.com/salsadigitalauorg/go-lagoon-log-forwarder.git
cd go-lagoon-log-forwarder

# Install dependencies and development tools
make deps
make install-tools

# Run quick development check
make quick
```

### Development with Makefile

We provide a comprehensive Makefile for all development tasks. To see all available commands:

```bash
make help
```

#### Quick Start Commands

```bash
make test           # Run all tests
make test-cover     # Run tests with coverage
make bench          # Run benchmarks
make lint           # Run linting
make ci             # Run full CI pipeline locally
```

#### Testing Commands

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make test-race` | Run tests with race detection |
| `make test-short` | Run tests in short mode (faster) |
| `make test-cover` | Run tests with coverage report |
| `make test-cover-html` | Generate HTML coverage report |
| `make bench` | Run benchmarks |
| `make bench-verbose` | Run benchmarks with verbose output |
| `make bench-compare` | Run benchmarks 5 times for comparison |

#### Code Quality Commands

| Command | Description |
|---------|-------------|
| `make lint` | Run golangci-lint |
| `make lint-fix` | Run linting with auto-fix |
| `make fmt` | Format code with gofmt + goimports |
| `make fmt-check` | Check if code is formatted correctly |
| `make vet` | Run go vet |
| `make security` | Run security scan with gosec |

#### Project Management

| Command | Description |
|---------|-------------|
| `make build` | Build the project |
| `make clean` | Clean build artifacts |
| `make deps` | Download and tidy dependencies |
| `make deps-update` | Update all dependencies |
| `make deps-verify` | Verify dependencies |
| `make install-tools` | Install development tools |

#### Workflow Commands

| Command | Description |
|---------|-------------|
| `make all` | Run fmt, vet, lint, and test |
| `make ci` | Run full CI pipeline locally |
| `make quick` | Quick development check (format, vet, test-short) |

#### Information Commands

| Command | Description |
|---------|-------------|
| `make version` | Show Go version and module info |
| `make info` | Show project information |

### Typical Development Workflow

```bash
# 1. Start development
make quick              # Quick check (format, vet, short tests)

# 2. Write code and tests

# 3. Run comprehensive tests
make test-cover         # Run tests with coverage
make bench             # Run benchmarks

# 4. Check code quality
make lint              # Run all linters
make security          # Security scan

# 5. Before committing
make ci                # Run full CI pipeline locally

# 6. Generate coverage report (optional)
make test-cover-html   # Open coverage.html in browser
```

### Manual Commands (without Makefile)

If you prefer not to use the Makefile:

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# View coverage report
go tool cover -html=coverage.out

# Format code
gofmt -s -w .
goimports -w .

# Run linting (requires golangci-lint)
golangci-lint run

# Run security scan (requires gosec)
gosec ./...
```

## Code Quality Standards

### Formatting

- Code must be formatted with `gofmt`
- Imports should be organized (use `goimports`)

### Testing

- All new code must have unit tests
- Maintain or improve test coverage
- Include benchmark tests for performance-critical code

### Documentation

- Public functions and types must have Go doc comments
- Update README.md for user-facing changes
- Include examples in doc comments where helpful

## Release Process

Releases are **fully automated**:

1. **PR Testing**: All PRs are automatically tested
2. **Merge**: When PR is merged to main, version is calculated from commit messages
3. **Release**: GitHub release is created with auto-generated changelog
4. **Tagging**: Git tag is created for the new version

### Version Calculation

- **No previous releases**: Starts at `v0.1.0`
- **Patch**: Bug fixes, documentation, refactoring
- **Minor**: New features, backward-compatible changes
- **Major**: Breaking changes (marked with `!` or `BREAKING CHANGE:`)

### Manual Release (if needed)

In rare cases, you can trigger a manual release:

```bash
# Create and push a tag
git tag v1.2.3
git push origin v1.2.3
```

## Getting Help

- **Issues**: Create a GitHub issue for bugs or feature requests
- **Discussions**: Use GitHub Discussions for questions
- **Code Review**: All PRs require review before merging

## License

By contributing, you agree that your contributions will be licensed under the same license as the project. 