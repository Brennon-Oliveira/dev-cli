# Command Structure Guide

This document explains how commands are structured in the Dev CLI project and provides guidelines for adding new commands.

## Overview

The Dev CLI uses the **Cobra** framework for command management. All commands are defined in the `cmd/` package and follow a consistent pattern to maintain code organization and testability.

## File Naming Convention

Commands are implemented in separate files following the pattern:
- `{command_name}.go` - Main command file (e.g., `run.go`, `shell.go`, `kill.go`)
- `root.go` - Entry point with the root command definition and global flags

## Command Structure Pattern

Every command follows this standard structure:

### 1. **Implementation Parameters Struct**

```go
type commandImplParams struct {
    args          []string
    pather        pather.Pather
    container     container.ContainerCLI
    devcontainer  devcontainer.DevContainerCLI
    // ... other dependencies
}
```

This struct holds all dependencies needed for command execution. It's designed for testability and dependency injection.

### 2. **Implementation Function**

```go
func commandImpl(p *commandImplParams) error {
    // Business logic here
    path := p.pather.GetPathFromArgs(p.args)
    absPath, _ := p.pather.GetAbsPath(path)
    
    return p.container.SomeOperation(absPath)
}
```

The implementation function contains the actual business logic. It:
- Is separated from the Cobra command definition
- Uses injected dependencies from the params struct
- Returns an error for proper error handling
- Can be easily tested with mock dependencies

### 3. **Cobra Command Definition**

```go
var commandCmd = &cobra.Command{
    Use:   "command [arguments]",
    Short: "Brief description",
    Long:  "Detailed description of what this command does",
    Args:  cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // Dependency initialization
        executor := exec.NewExecutor()
        pather := pather.NewPather(
            pather.WithExecutor(executor),
        )
        container := container.NewContainerCLI(
            container.WithExecutor(executor),
            container.WithPather(pather),
        )
        
        // Call implementation
        return commandImpl(&commandImplParams{
            args:      args,
            pather:    pather,
            container: container,
        })
    },
}
```

The Cobra command:
- Defines command metadata (Use, Short, Long, Examples)
- Handles argument validation (Args)
- Initializes all dependencies (using builder pattern)
- Calls the implementation function with injected params

### 4. **Command Registration**

```go
func init() {
    rootCmd.AddCommand(commandCmd)
}
```

All commands must be registered in the `init()` function to make them available to the root command.

## Key Patterns and Best Practices

### Dependency Injection

Commands use the **Builder Pattern** for dependency initialization:

```go
executor := exec.NewExecutor()
pather := pather.NewPather(
    pather.WithExecutor(executor),
)
devcontainer := devcontainer.NewDevContainerCLI(
    devcontainer.WithExecutor(executor),
)
```

Each dependency can be created with optional configuration functions.

### Argument Handling

Use Cobra's argument validation for ensuring correct usage:

- `cobra.MinimumNArgs(1)` - At least 1 argument required
- `cobra.MaximumNArgs(1)` - At most 1 argument allowed
- `cobra.ExactArgs(1)` - Exactly 1 argument required
- `cobra.NoArgs` - No arguments allowed

### Flags

Define command flags in the `init()` function:

```go
func init() {
    commandCmd.Flags().StringVarP(&someFlag, "flag-name", "f", "default", "Description")
    rootCmd.AddCommand(commandCmd)
}
```

For flags shared across commands, define them in `root.go`:

```go
var verboseFlag bool

func initLogger() {
    rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose output")
}
```

### Error Handling

Always use `RunE` instead of `Run`:

```go
RunE: func(cmd *cobra.Command, args []string) error {
    // Implementation
}
```

This allows proper error handling and reporting through Cobra.

### Path Normalization

Always use the `pather` package to handle paths:

```go
path := p.pather.GetPathFromArgs(p.args)      // Get path from arguments (defaults to ".")
absPath, _ := p.pather.GetAbsPath(path)       // Convert to absolute path
```

This ensures WSL compatibility and consistent path handling.

### Logging

Use the `logger` package for all output:

```go
logger.Info("Operation started")
logger.Verbose("Verbose information: %s", detail)
logger.Error("An error occurred: %s", err.Error())
```

## Command Examples Reference

### Simple Command Example: `up`

The `up` command is a minimal example:

```go
type upImplParams struct {
    args         []string
    pather       pather.Pather
    devcontainer devcontainer.DevContainerCLI
}

func upImpl(p *upImplParams) error {
    logger.Info("Iniciando projeto")
    path := p.pather.GetPathFromArgs(p.args)
    absPath, _ := p.pather.GetAbsPath(path)
    
    logger.Verbose("Rodando projeto na pasta %s", absPath)
    
    return p.devcontainer.Up(absPath)
}
```

### Complex Command Example: `run`

The `run` command shows a more complex example with multiple dependencies:

```go
type runImplParams struct {
    args         []string
    pather       pather.Pather
    devcontainer devcontainer.DevContainerCLI
    vscode       vscode.VSCode
}

func runImpl(p *runImplParams) error {
    logger.Info("Iniciando projeto")
    path := p.pather.GetPathFromArgs(p.args)
    absPath, _ := p.pather.GetAbsPath(path)
    
    if err := p.devcontainer.Up(absPath); err != nil {
        return err
    }
    
    workspaceURI, err := p.vscode.GetContainerWorkspaceURI(absPath)
    if err != nil {
        return err
    }
    
    return p.vscode.OpenWorkspaceByURI(workspaceURI)
}
```

### Command with Flags Example: `exec`

The `exec` command shows flag handling:

```go
var execPath string

var execCmd = &cobra.Command{
    Use:                "exec \"[comando]\"",
    Short:              "Executa um comando específico dentro do container",
    Args:               cobra.MinimumNArgs(1),
    DisableFlagParsing: false,
    RunE: func(cmd *cobra.Command, args []string) error {
        // ... initialization
    },
}

func init() {
    execCmd.Flags().StringVarP(&execPath, "path", "p", "", "Caminho do projeto (padrão '.')")
    execCmd.Flags().SetInterspersed(false)
    rootCmd.AddCommand(execCmd)
}
```

### Command with Custom Validation Example: `config`

The `config` command shows custom argument validation:

```go
var configCmd = &cobra.Command{
    Use:   "config [chave] [valor]",
    Short: "Gerencia as configurações da CLI",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) < 1 || len(args) > 2 {
            return fmt.Errorf("invalid arguments")
        }
        return nil
    },
    ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
        // Shell completion logic
    },
}
```

## Adding a New Command

To add a new command, follow these steps:

1. **Create the command file** (`cmd/new_command.go`)

2. **Define the implementation parameters**:
   ```go
   type newCommandImplParams struct {
       args []string
       // Add dependencies needed
   }
   ```

3. **Implement the business logic**:
   ```go
   func newCommandImpl(p *newCommandImplParams) error {
       // Implementation
   }
   ```

4. **Define the Cobra command**:
   ```go
   var newCommandCmd = &cobra.Command{
       Use:   "newcommand [args]",
       Short: "Brief description",
       Long:  "Detailed description",
       Args:  cobra.MaximumNArgs(1),
       RunE: func(cmd *cobra.Command, args []string) error {
           // Initialize dependencies
           // Call implementation
       },
   }
   ```

5. **Register the command**:
   ```go
   func init() {
       rootCmd.AddCommand(newCommandCmd)
   }
   ```

6. **Test the implementation** using mock dependencies (see internal package tests for examples)

## Available Dependencies

The following internal packages provide dependencies for commands:

- **`pather`** - Path resolution and normalization (WSL support)
- **`container`** - Docker/Podman container operations
- **`devcontainer`** - Dev Container CLI interface
- **`vscode`** - VS Code integration
- **`exec`** - System command execution
- **`config`** - Configuration management
- **`logger`** - Structured logging

Each dependency can be created with builder pattern options for customization and testing.

## Testing Commands

Commands are tested by:

1. Mocking dependencies using the auto-generated mocks in `internal/{package}/{interface}_mocks.go`
2. Creating test params structs with mock dependencies
3. Calling the implementation function directly
4. Asserting on results and mock call expectations

Example:
```go
mockContainer := new(container.MockContainerCLI)
mockContainer.On("SomeOperation", "path").Return(nil)

params := &commandImplParams{
    container: mockContainer,
}

err := commandImpl(params)

mockContainer.AssertExpectations(t)
```

## Common Patterns to Follow

✅ **Do:**
- Use dependency injection through params structs
- Separate implementation from Cobra command definition
- Return errors instead of calling logger.Fatal
- Use the pather package for all path operations
- Use the logger package for all output
- Initialize dependencies using the builder pattern

❌ **Don't:**
- Put business logic directly in the RunE function
- Call docker/podman commands directly with hardcoded strings
- Mix output methods (don't use fmt.Println and logger.Info)
- Assume absolute paths without using pather
- Initialize dependencies without options pattern

## Related Files

- `root.go` - Root command and global flags
- `internal/exec/executor.go` - System command execution interface
- `internal/pather/pather.go` - Path resolution interface
- `internal/container/container.go` - Container operations interface
- `internal/devcontainer/devcontainer.go` - Dev Container CLI interface
