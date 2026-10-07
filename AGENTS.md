# AGENTS.md — AI Agent Instructions for migtools/kopia

## Project Overview
This is the OADP fork of [Kopia](https://kopia.io), a fast and secure open-source backup/restore tool. Kopia provides encrypted, deduplicated, and compressed backups to cloud, network, or local storage. This fork (`migtools/kopia`) contains OADP-specific patches on the `oadp-dev` branch used by the OpenShift API for Data Protection (OADP) operator.

- **Primary Language**: Go
- **Module**: `github.com/kopia/kopia`
- **Default Branch**: `oadp-dev`

## Build Instructions
```bash
# Build the kopia binary
make install

# Build without HTML UI
make install-noui

# Build with race detector
make install-race
```

## Test Instructions
```bash
# Run all tests
go test ./...

# Run a specific test
go test ./path/to/package -run TestName

# Vet code
make vet
# or directly:
go vet ./...
```

## Linting
```bash
# Run linter (requires golangci-lint)
make lint

# Run linter with auto-fix
make lint-fix

# Platform-specific linting
make lint-linux
make lint-darwin
make lint-windows
make lint-all
```

Configuration: `.golangci.yml`

## Code Conventions
- Standard Go project layout with `internal/` for private packages
- CLI commands defined in `cli/` directory
- Repository/storage backends in `repo/`
- Snapshot management in `snapshot/`
- File system abstractions in `fs/`
- Follow existing patterns for error handling (wrapped errors with context)
- Notification system in `notification/`

## Project Structure
```
cli/           - CLI command definitions
cmd/           - Additional command binaries (OADP-specific additions)
fs/            - Filesystem abstraction layer
internal/      - Private packages (crypto, cache, compression, etc.)
notification/  - Notification system
repo/          - Repository and storage backends
site/          - Documentation website
snapshot/      - Snapshot management and policies
tests/         - Integration and end-to-end tests
tools/         - Build tooling
app/           - Desktop application (Electron-based)
```

## CI/CD
- GitHub Actions workflows in `.github/workflows/`:
  - `lint.yml` — Linting checks
  - `tests.yml` — Test suite
- Reproduce CI locally:
  ```bash
  make ci-setup  # Install tools and dependencies
  make lint      # Run linters
  go test ./...  # Run tests
  ```

## Common Tasks

### Adding a new storage backend
1. Implement the `repo/blob.Storage` interface in `repo/blob/`
2. Register the provider in `repo/blob/providers/`
3. Add CLI flags in `cli/`
4. Add tests in the corresponding `_test.go` files

### Adding a new CLI command
1. Create command file in `cli/`
2. Register in the command tree
3. Add tests

### Updating OADP-specific patches
- OADP patches live on the `oadp-dev` branch
- Keep patches minimal and rebasing-friendly against upstream `kopia/kopia`
