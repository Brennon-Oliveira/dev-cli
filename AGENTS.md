# AGENTS.md — Dev CLI Repository Guide

## Project Overview
**Dev CLI** is a Go CLI tool for managing Dev Container lifecycle and VS Code integration. It's a single-binary, multi-platform tool (Linux/macOS/Windows) with Cobra command structure.

**Key fact:** This is *not* a monorepo or web service—it's a focused CLI with command handlers in `cmd/`, reusable logic in `internal/`, and straightforward testing.

## Building & Testing

### Build the binary
```bash
go build -o dev main.go
```

**Important:** The CLI sets up the `PATH` variable from the shell on non-Windows systems (main.go:13-18). This is non-obvious and required for correct command execution in the container context.

### Run tests
```bash
go test -v ./...
```

- Tests use **testify** (`assert`, `require`) for assertions
- Mockery generates mocks in `_mocks.go` files (`.mockery.yml` configures: all packages auto-generate mocks in `internal/{pkg}/`)
- No integration test setup or database fixtures needed
- Tests run in isolation; safe to run in parallel

### Single package or test
```bash
go test -v ./internal/config
go test -v -run TestLoad ./internal/config
```

## Project Structure

- **`cmd/`** — Command definitions (Cobra). One file per command (`run.go`, `up.go`, `shell.go`, etc.)
  - `root.go` — Entry point with version injection and global `--verbose` flag
- **`internal/`** — Core logic, organized by concern:
  - `config/` — Config file loading/saving (JSON, `~/.dev/config`)
  - `container/` — Docker/Podman operations via shell `exec`
  - `devcontainer/` — Parse and resolve `devcontainer.json`
  - `vscode/` — Launch VS Code via `code` CLI
  - `exec/` — Execute commands in the container
  - `pather/` — Path translation for WSL compatibility (regex-based)
  - `logger/` — Structured logging with verbose mode
  - `env/` — Environment variable handling
- **`main.go`** — Tiny bootstrap; calls `cmd.Execute()`

## Key Quirks & Constraints

### Version handling
- Version is injected at build time via ldflags: `-ldflags="-X 'github.com/Brennon-Oliveira/dev-cli/cmd.Version=${{ github.ref_name }}'"`
- Default version is `"dev"` (set in `cmd/root.go:10`)
- Test builds will show version `"dev"` unless ldflags are explicitly set

### PATH issue (non-Windows)
The binary **re-evaluates the shell's PATH on startup** (main.go:13-18) to ensure container commands are found. This is a quirk agents might not expect—if a command seems to "not exist," the PATH handling may be the reason.

### Mock generation
- Run `mockery` separately if adding new interfaces to packages listed in `.mockery.yml`
- Mockery is configured to generate mocks in `internal/{pkg}/` as `{interface}_mocks.go` files
- To regenerate all mocks: `mockery --all` (requires mockery CLI installed)

### Docker/Podman abstraction
- The tool supports both Docker and Podman via config (`core.tool` key)
- Commands are run generically via `exec.Executor` interface; actual tool is selected at runtime
- See `.mockery.yml` for mocked interfaces: `config.Manager`, `exec.Executor`, `pather.Pather`, `devcontainer.Parser`, `container.Manager`, `logger.Logger`, `vscode.Handler`

## Testing Patterns

- Use `bytes.Buffer` or `io.Writer` for capturing output (logger tests: logger_test.go)
- Mock interfaces via Mockery-generated mocks
- Tests output is in Portuguese (matching the CLI); don't be surprised by colored terminal output in test logs

## CI/CD

Builds trigger on **version tags** (`v*`):
- `.github/workflows/build.yml` — Builds binaries for all platforms and creates artifacts
- `.github/workflows/snap-publish.yml` — Builds Snap package on `amd64`
- Version is read from the git tag (e.g., `v1.2.3`)

**Release checklist:** Tag the commit with `git tag v{version}` and push—GitHub Actions builds and uploads artifacts automatically.

## Common Commands for Development

```bash
# Rebuild after changes
go build -o dev main.go

# Run with verbose logging
./dev --verbose run /path/to/project

# Test a single command
./dev list
./dev up .
./dev kill .

# Format code (Go standard)
go fmt ./...
```

## Notes for Future Sessions

1. **Output is in Portuguese** — Don't assume English error messages. Check `internal/logger/` and command files for localized strings.
2. **Shell interaction is async** — Commands like `dev shell` inject an interactive shell; be careful with test automation of interactive flows.
3. **Path resolution is critical** — The `devcontainer.json` parsing in `internal/devcontainer/` and path translation in `internal/pather/` are the most complex pieces. Changes there need extra testing.
4. **Config persistence** — Config is stored in `~/.dev/config` as JSON. Changes to config keys or structure must update both `internal/config/` and any tests that mock the config manager.
5. **No external dependencies needed** — CLI is self-contained after `go build`. No Node, Python, or other runtimes required to develop or test.
