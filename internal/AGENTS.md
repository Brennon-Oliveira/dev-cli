# Internal Package Structure Guide

This document explains the architecture and organization of the `internal/` packages in the Dev CLI project. These packages contain all the core business logic that powers the CLI commands.

## Directory Structure

```
internal/
├── completer/          # Shell completion generation
├── config/             # Configuration management
├── constants/          # Application constants (currently empty)
├── container/          # Container operations (Docker/Podman)
│   └── container_utils/    # Container utility functions
├── devcontainer/       # Dev Container specification parsing
├── env/                # Environment variable handling
├── exec/               # System command execution
├── logger/             # Structured logging
│   └── logger_utils/   # Logging utility functions
├── pather/             # Path resolution and normalization
├── update/             # Application update logic
└── vscode/             # VS Code integration
```

## Core Patterns

All internal packages follow consistent patterns for maintainability and testability:

### Standard File Organization

Each package typically contains:

- **`{package}.go`** - Interface definitions
- **`{package}_builder.go`** - Builder pattern implementation with option functions
- **`{package}_impl.go`** - Concrete implementation of the interface
- **`{package}_mocks.go`** - Auto-generated mocks (via mockery)
- **`{package}_test.go`** - Unit tests
- **`{package}_utils.go`** - Utility functions (optional)
- **`*_utils/`** - Sub-package for complex utilities (optional)

### Builder Pattern

All packages use the builder pattern for initialization:

```go
// Interface
type SomeService interface {
    DoSomething() error
}

// Implementation with builder
type realService struct {
    dependency1 Dependency1
    dependency2 Dependency2
}

type Option func(*realService)

func NewService(opts ...Option) *realService {
    s := &realService{
        dependency1: defaultDep1,
        dependency2: defaultDep2,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

func WithDependency1(dep Dependency1) Option {
    return func(s *realService) {
        s.dependency1 = dep
    }
}
```

## Package Details

### 1. **exec** - System Command Execution

**Purpose:** Abstraction layer for executing system commands (sh, docker, devcontainer).

**Files:**
- `executor.go` - Interface definition
- `executor_builder.go` - Builder and options
- `executor_impl.go` - Main implementation
- `executor_posix.go` - POSIX-specific implementation (Linux/macOS)
- `executor_windows.go` - Windows-specific implementation
- `executor_mocks.go` - Auto-generated mocks
- `executor_test.go` - Unit tests

**Key Interface:**
```go
type Executor interface {
    Run(name string, args ...string) error
    RunWithOutput(output io.Writer, name string, args ...string) error
    RunInteractive(name string, args ...string) error
    RunDetached(name string, args ...string) error
    Output(name string, args ...string) ([]byte, error)
    CombinedOutput(name string, args ...string) (string, error)
}
```

**Methods:**
- `Run()` - Execute command and return error (output to stdout/stderr)
- `RunWithOutput()` - Execute with custom output writer
- `RunInteractive()` - Execute with TTY attachment (for interactive shells)
- `RunDetached()` - Execute in background
- `Output()` - Capture stdout only
- `CombinedOutput()` - Capture stdout and stderr combined

**Key Detail:** Handles platform-specific differences (Windows vs POSIX) using build tags and conditional compilation.

---

### 2. **pather** - Path Resolution and Normalization

**Purpose:** Handles all path operations with special support for WSL (Windows Subsystem for Linux).

**Files:**
- `pather.go` - Interface definition
- `pather_builder.go` - Builder and options
- `pather_impl.go` - Implementation with WSL detection
- `pather_mocks.go` - Auto-generated mocks
- `pather_test.go` - Unit tests

**Key Interface:**
```go
type Pather interface {
    GetPathFromArgs(args []string) string
    GetAbsPath(path string) (string, error)
    GetHostPath(absPath string) (string, error)
    IsInWSL() bool
}
```

**Methods:**
- `GetPathFromArgs()` - Extract path from command arguments (defaults to ".")
- `GetAbsPath()` - Convert to absolute path
- `GetHostPath()` - Convert to Windows path (when in WSL)
- `IsInWSL()` - Detect if running in WSL

**Key Detail:** Uses `wslpath` command to convert paths between WSL and Windows when needed. This is critical for VS Code integration on Windows hosts.

---

### 3. **config** - Configuration Management

**Purpose:** Manages persistent CLI configuration stored in `~/.dev/config` (JSON format).

**Files:**
- `config.go` - Interface definition
- `config_builder.go` - Builder and options
- `config_impl.go` - Implementation
- `config_default.go` - Default configuration values
- `config_mocks.go` - Auto-generated mocks
- `config_test.go` - Unit tests
- `config_utils.go` - Utility functions
- `global_config.go` - Global config structure

**Key Interface:**
```go
type Config interface {
    GetConfigPath() (string, error)
    HasConfigFile() bool
    Load() GlobalConfig
    LoadByKey(key string) string
    TrySave(key string, value string) (string, error)
    Save(key string, value string) error
    InterativeSelect(key string) (string, error)
    ValidateKey(key string) bool
}
```

**Configuration Structure:**
```go
type GlobalConfig struct {
    Core CoreConfig `json:"core"`
}

type CoreConfig struct {
    Tool string `json:"tool"` // "docker" or "podman"
}
```

**Methods:**
- `Load()` - Load entire config from file
- `LoadByKey()` - Get a specific config value
- `TrySave()` - Attempt to save and validate configuration
- `Save()` - Persist config to file
- `ValidateKey()` - Check if key is valid
- `InterativeSelect()` - Interactive menu for user selection

**Key Detail:** Configuration handlers define valid values and display options for each setting.

---

### 4. **devcontainer** - Dev Container CLI Integration

**Purpose:** Interface with the Dev Container CLI (`devcontainer` command) and parse `devcontainer.json` specifications.

**Files:**
- `devcontainer.go` - Interface and configuration structures
- `devcontainer_builder.go` - Builder and options
- `devcontainer_impl.go` - Implementation
- `devcontainer_mocks.go` - Auto-generated mocks
- `devcontainer_test.go` - Unit tests

**Key Interface:**
```go
type DevContainerCLI interface {
    Up(workspace string) error
    GetWorkspaceFolder(absPath string) (string, error)
    ReadConfiguration(absPath string) (*DevContainerConfiguration, error)
    RunInteractive(path string, command string) error
    OpenShell(path string) error
}
```

**Configuration Structure:**
```go
type DevContainerConfiguration struct {
    Workspace DevContainerConfiguration_Workspace `json:"workspace"`
}

type DevContainerConfiguration_Workspace struct {
    WorkspaceFolder string `json:"workspaceFolder"`
}
```

**Methods:**
- `Up()` - Build and start container
- `GetWorkspaceFolder()` - Resolve workspace folder from devcontainer.json
- `ReadConfiguration()` - Parse devcontainer.json
- `RunInteractive()` - Execute command with output
- `OpenShell()` - Start interactive shell session

**Key Detail:** Reads and interprets the `devcontainer.json` configuration file to determine container settings, workspace location, and mounted volumes.

---

### 5. **container** - Container Operations

**Purpose:** Low-level container management using Docker or Podman.

**Files:**
- `container.go` - Interface definition
- `container_builder.go` - Builder and options
- `container_impl.go` - Implementation
- `container_mocks.go` - Auto-generated mocks
- `container_test.go` - Unit tests
- `container_utils/` - Utility functions for container operations

**Key Interface:**
```go
type ContainerCLI interface {
    KillContainer(absPath string) error
    GetAllContainers() (*container.ContainerListOptions, error)
    GetContainerLogs(containerID string) (string, error)
    GetContainerPorts(containerID string) ([]string, error)
    // ... additional methods
}
```

**Methods:**
- `KillContainer()` - Stop and remove container and related services
- `GetAllContainers()` - List running containers
- `GetContainerLogs()` - Retrieve container logs
- `GetContainerPorts()` - Get port mappings

**Key Detail:** Uses the configured tool (Docker or Podman) based on user preferences. Automatically discovers related containers (e.g., compose services) and manages them together.

---

### 6. **vscode** - VS Code Integration

**Purpose:** Launch and manage VS Code connections to containers.

**Files:**
- `vscode.go` - Interface definition
- `vscode_builder.go` - Builder and options
- `vscode_impl.go` - Implementation
- `vscode_mocks.go` - Auto-generated mocks
- `vscode_test.go` - Unit tests

**Key Interface:**
```go
type VSCode interface {
    OpenWorkspaceByURI(workspaceURI string) error
    GetContainerWorkspaceURI(absPath string) (string, error)
}
```

**Methods:**
- `OpenWorkspaceByURI()` - Launch VS Code with remote container URI
- `GetContainerWorkspaceURI()` - Generate vscode-remote:// URI

**Key Detail:** Handles WSL path conversion to ensure Windows VS Code instances can properly connect to containers in WSL environments.

---

### 7. **logger** - Structured Logging

**Purpose:** Centralized logging with verbose mode support.

**Files:**
- `logger.go` - Interface and configuration
- `logger_builder.go` - Builder and options
- `logger_impl.go` - Implementation
- `logger_writer.go` - Output writer management
- `logger_utils/` - Formatting utilities
- `logger_mocks.go` - Auto-generated mocks
- `logger_test.go` - Unit tests

**Key Interface:**
```go
type Logger interface {
    Info(format string, args ...interface{})
    Verbose(format string, args ...interface{})
    Error(format string, args ...interface{})
    Success(format string, args ...interface{})
}
```

**Methods:**
- `Info()` - General information messages
- `Verbose()` - Debug/verbose messages (only shown with --verbose flag)
- `Error()` - Error messages
- `Success()` - Success/completion messages

**Key Detail:** Output is in Portuguese. Uses ANSI color codes for terminal output. Verbose mode is controlled by global `--verbose` flag in root command.

---

### 8. **env** - Environment Variables

**Purpose:** Safe environment variable lookup with defaults.

**Files:**
- `lookup_env.go` - Lookup functions

**Key Functions:**
- `LookupEnv()` - Safely retrieve environment variables

---

### 9. **update** - Application Updates

**Purpose:** Check for and apply CLI updates.

**Files:**
- Update-related files for version checking and self-updates

---

### 10. **completer** - Shell Completion

**Purpose:** Generate shell completion scripts (bash, zsh, fish).

**Files:**
- Completion generation logic for different shell types

---

## Dependency Graph

```
cmd/ (commands)
  ├── exec.Executor
  ├── pather.Pather
  ├── devcontainer.DevContainerCLI
  ├── container.ContainerCLI
  ├── vscode.VSCode
  ├── config.Config
  └── logger.Logger

exec.Executor (no dependencies)

pather.Pather
  └── exec.Executor

devcontainer.DevContainerCLI
  └── exec.Executor

container.ContainerCLI
  ├── exec.Executor
  ├── pather.Pather
  └── config.Config

vscode.VSCode
  ├── exec.Executor
  ├── pather.Pather
  └── devcontainer.DevContainerCLI

config.Config (no dependencies)

logger.Logger (no dependencies)
```

## Testing Approach

All packages are tested using:

1. **Mock Generation** - Mockery auto-generates mocks from interfaces
2. **Dependency Injection** - Tests inject mock dependencies via builder options
3. **Unit Tests** - Each function/method has isolated tests
4. **Table-Driven Tests** - Complex scenarios use table-driven test patterns

Example test pattern:
```go
func TestSomeOperation(t *testing.T) {
    mockExec := new(exec.MockExecutor)
    mockExec.On("Run", "command", "arg").Return(nil)
    
    service := NewService(
        WithExecutor(mockExec),
    )
    
    err := service.SomeOperation()
    
    require.NoError(t, err)
    mockExec.AssertExpectations(t)
}
```

## Adding New Internal Packages

When creating a new internal package:

1. **Define the interface** in `{package}.go`
2. **Implement builder** in `{package}_builder.go`
3. **Implement logic** in `{package}_impl.go`
4. **Write tests** in `{package}_test.go`
5. **Regenerate mocks** using `mockery` CLI
6. **Add to .mockery.yml** configuration if needed

## Key Principles

✅ **Do:**
- Use interfaces for dependencies (enables mocking)
- Use builder pattern with options
- Separate concerns into focused packages
- Return errors instead of calling logger.Fatal
- Handle both Docker and Podman via config
- Support WSL path conversion
- Test with mock dependencies

❌ **Don't:**
- Mix business logic with system calls
- Hardcode Docker/Podman in implementations
- Assume absolute paths without normalization
- Put test-only code in main implementations
- Create circular dependencies between packages
- Ignore WSL compatibility

## Constants Package

The `constants/` package is currently empty and reserved for future global constants. When adding constants:

1. Define them in `internal/constants/constants.go`
2. Import as `constants.CONSTANT_NAME`
3. Use for magic strings and configuration values

## Related Documentation

- `cmd/AGENTS.md` - Command structure guide
- `AGENTS.md` (root) - General project guide
- `docs/` - User-facing documentation
