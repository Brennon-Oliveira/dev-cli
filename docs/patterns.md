# Patterns and Conventions

To maintain consistency and code maintainability, follow these patterns:

## Code Organization

- **Encapsulation** - Logic that doesn't need external exposure should be in `internal/`.
- **Container Engine Abstraction** - Never call `docker` or `podman` directly with hardcoded strings. Always use `config.Load().Core.Tool` to respect user preferences.
- **System Calls** - Use abstraction functions in `internal/exec/` to ensure processes are handled correctly across platforms.

## Error Handling

- Return errors instead of calling `log.Fatal` or `panic()` (except at entry point in `main.go` or `root.go`).
- Wrap errors with context when needed: `fmt.Errorf("operation context: %w", err)`
- Always propagate errors up the call stack for proper handling at the command level.

### Error Example

```go
func SomeOperation() error {
    if err := someCall(); err != nil {
        return fmt.Errorf("failed to complete operation: %w", err)
    }
    return nil
}
```

## Cross-Platform Compatibility

The project supports Windows, Linux, macOS, and WSL.

### Platform-Specific Code

Use file suffixes for platform-specific implementations:
- `executor_posix.go` - Linux/macOS specific
- `executor_windows.go` - Windows specific

### Path Handling

Always use `filepath.Join()` for path construction:

```go
// Correct
path := filepath.Join(homeDir, ".dev", "config")

// Incorrect
path := homeDir + "/.dev/config"  // Won't work on Windows
```

### Architecture Support

Code runs on multiple architectures:
- amd64 (x86_64)
- arm64

Avoid architecture-specific assumptions.

## Dependency Injection Pattern

All packages use builder pattern with functional options:

```go
type Option func(*implementation)

func NewService(opts ...Option) *implementation {
    s := &implementation{
        // default values
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

func WithDependency(dep SomeDependency) Option {
    return func(s *implementation) {
        s.dependency = dep
    }
}
```

Benefits:
- Easy to mock for testing
- Optional configuration
- Clear dependencies
- Backward compatible

## Interface-Based Design

Create interfaces for external dependencies:

```go
type FileReader interface {
    ReadFile(path string) ([]byte, error)
}

type MyService struct {
    reader FileReader  // Dependency injection
}

func NewMyService(reader FileReader) *MyService {
    return &MyService{reader: reader}
}
```

This allows:
- Multiple implementations
- Easy mocking in tests
- Clear contracts

## Logging

Use the `logger` package for all output:

```go
import "github.com/Brennon-Oliveira/dev-cli/internal/logger"

logger.Info("Starting operation")
logger.Verbose("Debug info shown with --verbose")
logger.Error("An error occurred")
logger.Success("Operation completed")
```

**Note:** Output is in Portuguese. Use color-coded messages for better UX.

## Testing Patterns

### Mock Dependencies

Use auto-generated mocks from mockery:

```go
mockExec := new(exec.MockExecutor)
mockExec.On("Run", "command", "arg").Return(nil)

// Use in test
```

### Table-Driven Tests

For complex scenarios:

```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"case1", "input1", "output1", false},
    {"case2", "input2", "output2", true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := Function(tt.input)
        if (err != nil) != tt.wantErr {
            t.Errorf("unexpected error")
        }
        if got != tt.want {
            t.Errorf("got %q, want %q", got, tt.want)
        }
    })
}
```

### Assertion Library

Use `testify` for cleaner assertions:

```go
import "github.com/stretchr/testify/require"

require.NoError(t, err)
require.Equal(t, expected, actual)
require.NotNil(t, result)
```

## Configuration Extension

To add new configuration options:

1. **Update GlobalConfig struct**:
   ```go
   type GlobalConfig struct {
       Core CoreConfig `json:"core"`
       // Add new field
   }
   ```

2. **Add configuration handler** in `internal/config/config_utils.go`:
   ```go
   handlers["new.option"] = ConfigHandler{
       ValidValues: []string{"value1", "value2"},
       Display:     "Display name",
   }
   ```

3. **Write tests** to verify persistence and validation

## Naming Conventions

### Functions

- **Actions**: Use verbs (Get, Set, Create, Delete, Run, Execute)
- **Helpers**: Use prefixes (IsValid, HasValue, TryParse)

```go
// Good
func GetAbsPath(path string) (string, error)
func RunCommand(cmd string) error
func IsValidPath(path string) bool

// Bad
func Path(p string) (string, error)
func Command(cmd string) error
func Valid(p string) bool
```

### Types

- **Interfaces**: Use `-er` suffix or descriptive names
- **Implementations**: Use concrete names

```go
// Good
type Reader interface { ... }
type realReader struct { ... }

type ContainerCLI interface { ... }
type containerImpl struct { ... }

// Bad
type IReader interface { ... }
type Reader struct { ... }
```

### Variables

- Use descriptive names
- Avoid single letters except for loops/iterators

```go
// Good
absPath, err := pather.GetAbsPath(path)
for i, item := range items { ... }

// Bad
p, e := pather.GetAbsPath(path)
ap := absPath
```

## File Organization

Each package should be self-contained:

```
internal/mypackage/
├── mypackage.go           # Interface
├── mypackage_builder.go   # Constructor and options
├── mypackage_impl.go      # Implementation
├── mypackage_mocks.go     # Auto-generated mocks
├── mypackage_test.go      # Tests
├── mypackage_utils.go     # Helper functions (optional)
└── mypackage_utils/       # Utility subpackage (optional)
```

## Comments and Documentation

- **Public functions**: Start with function name
  ```go
  // GetAbsPath converts a relative path to absolute.
  func GetAbsPath(path string) (string, error) { ... }
  ```

- **Complex logic**: Explain the "why"
  ```go
  // IsInWSL checks for WSL-specific environment variables
  // to determine if running in Windows Subsystem for Linux.
  ```

- **Deprecated**: Mark clearly
  ```go
  // Deprecated: Use NewFunction instead.
  func OldFunction() { ... }
  ```

## Performance Considerations

- Avoid repeated path conversions (cache results)
- Minimize shell command executions
- Use interfaces to defer expensive operations
- Profile before optimizing

## Security Considerations

- Validate all user input
- Avoid `eval()` or dynamic code execution
- Handle sensitive paths (config, credentials) carefully
- Check file permissions before reading/writing

## Related Documentation

- `cmd/AGENTS.md` - Command implementation patterns
- `internal/AGENTS.md` - Internal package architecture
- `architecture.md` - System-level architecture
- `commands.md` - Command creation guide
