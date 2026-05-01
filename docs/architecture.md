# Dev CLI Architecture

Dev CLI is a Go-based tool designed to orchestrate the lifecycle of Dev Containers, integrating natively with VS Code while handling the complexities of environments like WSL (Windows Subsystem for Linux).

## Overview

The application follows a standard modular structure in Go:

- **`main.go`** - Entry point that initializes command execution
- **`cmd/`** - Command definitions using the **Cobra** framework. Each file represents a command or subcommand
- **`internal/`** - Protected business logic, inaccessible to external modules
  - **`internal/config/`** - Manages global CLI configuration (e.g., Docker vs Podman selection)
  - **`internal/container/`** - Core logic for interacting with container engines and Dev Container CLI
  - **`internal/exec/`** - System command execution abstraction
  - **`internal/pather/`** - Path resolution with WSL support
  - **`internal/devcontainer/`** - Dev Container specification parsing
  - **`internal/vscode/`** - VS Code integration
  - **`internal/logger/`** - Structured logging
  - **`internal/env/`** - Environment variable handling
  - **`internal/update/`** - Self-update functionality
  - **`internal/completer/`** - Shell completion generation

## Execution Flow

1. User executes a command (e.g., `dev run .`)
2. `main.go` calls `cmd.Execute()`
3. Cobra identifies and routes to the appropriate command handler (e.g., `cmd/run.go`)
4. Command validates arguments and initializes dependencies
5. Command calls internal package functions with injected dependencies
6. Internal packages execute system commands through the abstraction layers
7. Results are returned and formatted for user display

### Example: The `run` Command

```
User Input: dev run .
    ↓
main.go calls cmd.Execute()
    ↓
cmd/run.go RunE function
    ↓
Create dependencies: executor, pather, devcontainer, vscode
    ↓
runImpl() executes:
  1. pather.GetPathFromArgs() → resolve "." to absolute path
  2. devcontainer.Up() → build and start container
  3. vscode.GetContainerWorkspaceURI() → generate remote URI
  4. vscode.OpenWorkspaceByURI() → launch VS Code
    ↓
Result: VS Code opens connected to container
```

## Key Design Principles

### Dependency Injection

All components use the **builder pattern** with functional options:

```go
executor := exec.NewExecutor()
pather := pather.NewPather(
    pather.WithExecutor(executor),
)
```

This enables:
- Easy mocking for tests
- Flexible configuration
- Loose coupling between packages

### Interface-Based Design

Core logic is abstracted behind interfaces:

```go
type Executor interface {
    Run(name string, args ...string) error
    RunInteractive(name string, args ...string) error
    // ...
}
```

This allows:
- Multiple implementations (Windows vs POSIX)
- Mock implementations for testing
- Easy substitution without code changes

### Error Handling

Errors bubble up through the call stack:
- Internal packages return errors
- Command handlers convert errors to user-friendly messages
- No `panic()` or `os.Exit()` except at the entry point

## WSL Integration

One of Dev CLI's key differentiators is robust WSL support. The tool:

1. **Detects WSL environment** - Checks for WSL-specific environment variables
2. **Converts paths** - Uses `wslpath` to convert between WSL and Windows paths
3. **Manages URI schemes** - Ensures VS Code (running on Windows) can connect to containers
4. **Handles workspaces** - Correctly resolves workspace folders in both WSL and Windows contexts

The `pather` package handles all path conversion, ensuring transparent operation across platforms.

## Platform Support

- **Linux** - Native support
- **macOS** - Native support
- **Windows (WSL)** - Full support with path conversion
- **Windows (native)** - Direct container execution (via Docker Desktop)

Platform-specific code uses Go build tags:
- `executor_posix.go` - Linux/macOS specific
- `executor_windows.go` - Windows specific

## Configuration

Configuration is stored in `~/.dev/config` as JSON:

```json
{
  "core": {
    "tool": "docker"
  }
}
```

The `config` package provides:
- Persistent storage
- Validation
- Interactive selection menus
- Per-command configuration

## Dependencies

Dev CLI uses minimal external dependencies:

- **Cobra** - Command-line framework
- **Testify** - Testing assertions
- **Mockery** - Mock generation

No additional runtimes are required (no Node.js, Python, etc.).

## Execution Environment

On non-Windows systems, the CLI re-evaluates the shell's PATH on startup (`main.go:13-18`) to ensure container commands and development tools are discovered correctly. This is essential for finding:
- Docker/Podman CLI
- Dev Container CLI
- System utilities

## Related Documentation

- `cmd/AGENTS.md` - Command implementation patterns
- `internal/AGENTS.md` - Internal package structure
- `docs/commands.md` - Command creation guide
- `docs/patterns.md` - Development patterns
