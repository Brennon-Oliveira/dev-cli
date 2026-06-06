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

## Codex Harness

This repository includes a repo-local Codex harness in `.agents/` and `.codex/`. Harness files are written in English for agent clarity; keep existing Portuguese CLI output and tests unchanged unless a task explicitly touches user-facing copy.

### Repo skills

Codex can discover checked-in skills from `.agents/skills/`:

- `$release-version` - Use for version release preparation, release notes, verification, tag pushes, and GitHub Release create/update work.
- `$new-internal-package` - Use when adding a package under `internal/` with the interface/builder/implementation/test/mock pattern.
- `$fix-bug` - Use for reproducing a defect, applying a narrow patch, adding regression coverage, and running targeted plus full tests.
- `$add-command` - Use when adding a Cobra command under `cmd/` with dependency injection, registration, tests, and docs updates when behavior changes.
- `$commit-changes` - Use when splitting the current worktree into focused commits that strictly follow Conventional Commits 1.0.0.

See `.agents/README.md` for trigger examples.

### Hooks and rules

- `.codex/hooks.json` defines post-edit automation and a non-blocking Stop hook. The post-edit hook runs `git diff --check` for all file changes and `go fmt ./...` for Go changes without blocking. It runs `go mod tidy` for Go/module changes and package-level `go test -v` when `*_test.go` files change; failures are reported to the agent with a non-zero exit so they can be fixed before finishing. The Stop hook emits valid empty JSON because Codex validates Stop hook stdout as JSON.
- `.codex/rules/dev-cli.rules` allows common safe verification/read-only commands, including `go mod tidy`, and prompts before release-mutating commands such as `git push` and `gh release create/upload/edit`.
- To validate rules, run examples with:
  ```bash
  codex execpolicy check --pretty --rules .codex/rules/dev-cli.rules -- go test ./...
  ```

## Testing Patterns

- Use `bytes.Buffer` or `io.Writer` for capturing output (logger tests: logger_test.go)
- Mock interfaces via Mockery-generated mocks
- Tests output is in Portuguese (matching the CLI); don't be surprised by colored terminal output in test logs

## CI/CD

Builds trigger on **version tags** (`v*`):
- `.github/workflows/build.yml` — Builds binaries for all platforms and creates artifacts
- `.github/workflows/snap-publish.yml` — Builds Snap package on `amd64`
- Version is read from the git tag (e.g., `v1.2.3`)

**Release checklist:** Use `$release-version` when available. Tag the commit with `git tag v{version}` and push only after verification and release-note review. GitHub Actions builds artifacts and creates or updates the GitHub Release for the tag. Because Snap publishing runs on the GitHub `release: released` event, publish only after release notes and tag state are correct.

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
