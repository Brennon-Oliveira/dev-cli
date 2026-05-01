# Development and Contribution

## Prerequisites

- Go 1.25.6 or higher
- Docker or Podman installed
- Dev Container CLI (`npm install -g @devcontainers/cli`)
- VS Code installed (optional, for testing `run` and `open` commands)

## Local Configuration

The CLI stores configuration in `~/.dev/config` (JSON format). You can configure the container tool manually:

```json
{
  "core": {
    "tool": "docker"
  }
}
```

Valid values for `tool`:
- `"docker"` - Use Docker (default)
- `"podman"` - Use Podman

## Building the Project

To build the binary locally:

```bash
go build -o dev .
```

Then move it to your PATH or run directly:

```bash
./dev --help
```

To build with version information:

```bash
go build -ldflags="-X 'github.com/Brennon-Oliveira/dev-cli/cmd.Version=v1.0.0'" -o dev .
```

## Running Tests

Run all tests:

```bash
go test -v ./...
```

Run tests for a specific package:

```bash
go test -v ./internal/config
```

Run a specific test:

```bash
go test -v -run TestLoad ./internal/config
```

Run tests in parallel:

```bash
go test -v -parallel 4 ./...
```

## Formatting Code

Format all code to Go standards:

```bash
go fmt ./...
```

## Regenerating Mocks

If you add new interfaces to packages listed in `.mockery.yml`, regenerate mocks:

```bash
mockery --all
```

Mocks are automatically generated in `internal/{package}/{interface}_mocks.go` files.

## Project Structure for Development

```
dev-cli/
├── main.go              # Entry point
├── cmd/                 # Commands
│   ├── AGENTS.md        # Command patterns guide
│   ├── root.go          # Root command
│   ├── run.go           # Run command
│   └── ...
├── internal/            # Core logic
│   ├── AGENTS.md        # Package structure guide
│   ├── exec/            # Command execution
│   ├── pather/          # Path resolution
│   ├── config/          # Configuration
│   ├── container/       # Container operations
│   ├── devcontainer/    # Dev Container CLI
│   ├── vscode/          # VS Code integration
│   ├── logger/          # Logging
│   └── ...
├── docs/                # Documentation
│   ├── architecture.md
│   ├── commands.md
│   ├── development.md
│   └── patterns.md
└── README.md            # User documentation
```

## Development Workflow

1. **Create a feature branch**: `git checkout -b feature/my-feature`
2. **Make changes** in `cmd/` or `internal/`
3. **Write tests** in `*_test.go` files
4. **Run tests**: `go test -v ./...`
5. **Format code**: `go fmt ./...`
6. **Commit**: `git commit -m "descriptive message"`
7. **Push and create PR**

## Testing Considerations

### WSL Testing

If you modify path resolution or process execution logic, test in WSL:

```bash
# Inside WSL
./dev run .
./dev shell
./dev list
```

WSL has specific:
- Path mapping rules (Windows ↔ Linux paths)
- Network configuration
- File system behavior

### Docker vs Podman

Test with both container engines if modifying container operations:

```bash
dev config --global core.tool docker
./dev run .

dev config --global core.tool podman
./dev run .
```

### Platform Testing

If modifying platform-specific code, test on:
- Linux
- macOS
- Windows (WSL)
- Windows (native with Docker Desktop)

Platform-specific files use build tags:
- `executor_posix.go` - Linux/macOS
- `executor_windows.go` - Windows

## Common Development Tasks

### Adding a New Command

1. Create `cmd/mycommand.go` following the pattern in `cmd/AGENTS.md`
2. Write tests in the same file or `cmd/mycommand_test.go`
3. Run: `go build -o dev . && ./dev mycommand --help`

### Adding Configuration Options

1. Update `GlobalConfig` struct in `internal/config/global_config.go`
2. Add handler in `internal/config/config_utils.go`
3. Update tests in `internal/config/config_test.go`

### Modifying Path Resolution

1. Update logic in `internal/pather/pather_impl.go`
2. Test with WSL paths
3. Update tests in `internal/pather/pather_test.go`

### Adding Logging

Use the logger package:

```go
import "github.com/Brennon-Oliveira/dev-cli/internal/logger"

logger.Info("Message")
logger.Verbose("Debug info")
logger.Error("Error message")
```

Output is in Portuguese and includes color formatting.

## Debugging

### Enable Verbose Output

Run commands with `--verbose` flag:

```bash
./dev run --verbose .
```

This shows all executed commands and debug information.

### Print Debug Information

Use `logger.Verbose()`:

```go
logger.Verbose("Value of x: %v", x)
```

Only shown when `--verbose` flag is used.

## Release Process

Releases are triggered by git tags:

1. Tag the commit: `git tag v1.2.3`
2. Push the tag: `git push origin v1.2.3`
3. GitHub Actions automatically builds and uploads artifacts

Supported platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

Snap package is also built for Linux.

## Performance Considerations

- Container startup is the main bottleneck (handled by Docker/Podman)
- Path resolution on WSL adds minimal overhead
- Logger is designed for minimal performance impact

## Documentation

- Update `docs/` when changing functionality
- Update command descriptions in `cmd/{command}.go`
- Add examples to command Long descriptions
- Keep `README.md` synchronized

## Related Documentation

- `cmd/AGENTS.md` - Command structure patterns
- `internal/AGENTS.md` - Internal package structure
- `architecture.md` - System architecture
- `patterns.md` - Development patterns
