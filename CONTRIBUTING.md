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

### Setup

```bash
# Clone the repository
git clone https://github.com/salsadigitalauorg/go-lagoon-log-forwarder.git
cd go-lagoon-log-forwarder

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run linting (optional, requires golangci-lint)
golangci-lint run
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# View coverage report
go tool cover -html=coverage.out
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