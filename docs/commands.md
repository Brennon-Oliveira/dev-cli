# Creating New Commands

Dev CLI uses the [Cobra](https://github.com/spf13/cobra) framework for command management.

## Command Structure

To create a new command, add a `.go` file in the `cmd/` directory. For example, `cmd/example.go`:

```go
package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/spf13/cobra"
)

// Implementation parameters
type exampleImplParams struct {
	args []string
	// Add dependencies here
}

// Implementation function
func exampleImpl(p *exampleImplParams) error {
	logger.Info("Running example command")
	if len(p.args) > 0 {
		logger.Info("Argument: %s", p.args[0])
	}
	return nil
}

// Cobra command definition
var exampleCmd = &cobra.Command{
	Use:   "example [argument]",
	Short: "Brief description of the command",
	Long:  "Detailed description of what this command does.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize dependencies
		executor := exec.NewExecutor()
		
		return exampleImpl(&exampleImplParams{
			args: args,
		})
	},
}

// Register the command
func init() {
	rootCmd.AddCommand(exampleCmd)
}
```

## Command Patterns

### 1. Argument Validation

Use Cobra's argument validators:

```go
Args: cobra.MinimumNArgs(1),     // At least 1 argument
Args: cobra.MaximumNArgs(1),     // At most 1 argument
Args: cobra.ExactArgs(1),        // Exactly 1 argument
Args: cobra.NoArgs,              // No arguments allowed
```

### 2. Use RunE Instead of Run

Always prefer `RunE`:

```go
RunE: func(cmd *cobra.Command, args []string) error {
	// Implementation
	return nil  // Return errors directly
}
```

This enables proper error handling and reporting through Cobra.

### 3. Separate Logic from Command Definition

Keep business logic in a separate function that accepts injected dependencies:

```go
type myCommandImplParams struct {
	args      []string
	executor  exec.Executor
	pather    pather.Pather
}

func myCommandImpl(p *myCommandImplParams) error {
	// Business logic here
}
```

This pattern:
- Makes testing easier (inject mocks)
- Keeps command definitions clean
- Enables code reuse

### 4. Path Normalization

Always normalize paths using the `pather` package:

```go
path := p.pather.GetPathFromArgs(p.args)        // Get path from args (defaults to ".")
absPath, _ := p.pather.GetAbsPath(path)         // Convert to absolute path
```

This ensures WSL compatibility and consistent handling across platforms.

### 5. Use Logger for Output

Use the `logger` package instead of `fmt.Println`:

```go
logger.Info("Normal output")
logger.Verbose("Debug information")
logger.Error("Error message")
logger.Success("Success message")
```

### 6. Command Flags

Define flags in the `init()` function:

```go
var myFlag string

func init() {
	exampleCmd.Flags().StringVarP(&myFlag, "flag-name", "f", "default", "Flag description")
	rootCmd.AddCommand(exampleCmd)
}
```

For global flags shared by all commands, define in `root.go`.

## Best Practices

✅ **Do:**
- Use dependency injection through params structs
- Separate implementation from command definition
- Return errors instead of using logger.Fatal
- Use pather for all path operations
- Use logger for all output
- Initialize dependencies using builder pattern
- Write tests for implementation functions

❌ **Don't:**
- Put business logic directly in RunE
- Call docker/podman directly with hardcoded strings
- Mix output methods (don't combine fmt and logger)
- Assume absolute paths without using pather
- Initialize dependencies without builder options

## Testing Your Command

Test the implementation function directly with mock dependencies:

```go
import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestExampleImpl(t *testing.T) {
	mockExec := new(exec.MockExecutor)
	mockExec.On("Run", "command").Return(nil)
	
	params := &exampleImplParams{
		args:     []string{"test"},
		executor: mockExec,
	}
	
	err := exampleImpl(params)
	
	require.NoError(t, err)
	mockExec.AssertExpectations(t)
}
```

## Available Packages

Import from `internal/` for these capabilities:

- **`exec`** - System command execution
- **`pather`** - Path resolution
- **`container`** - Container operations
- **`devcontainer`** - Dev Container CLI
- **`vscode`** - VS Code integration
- **`config`** - Configuration management
- **`logger`** - Logging
- **`env`** - Environment variables

## Full Example: List Containers Command

See the actual commands in `cmd/` directory for complete examples:
- `cmd/run.go` - Complex command with multiple dependencies
- `cmd/up.go` - Simple command
- `cmd/shell.go` - Interactive command
- `cmd/config.go` - Command with custom validation

Refer to `cmd/AGENTS.md` for detailed patterns and examples.

## Related Documentation

- `cmd/AGENTS.md` - Comprehensive command structure guide
- `internal/AGENTS.md` - Internal package documentation
- `architecture.md` - Overall system architecture
