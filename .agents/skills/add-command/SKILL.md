---
name: add-command
description: Use when adding a new Cobra command to Dev CLI, including cmd/AGENTS.md conventions, dependency injection params, command registration, tests, and docs updates when user-facing behavior changes.
---

# Add Command

Use this workflow for new Dev CLI commands.

## Command Design

1. Read `cmd/AGENTS.md` before editing command files.
2. Create `cmd/{command}.go` and follow the existing pattern:
   - `{command}ImplParams` for injected dependencies and args.
   - `{command}Impl` for testable behavior.
   - A `cobra.Command` with `Use`, `Short`, `Long`, `Args`, and `RunE`.
   - `init()` registration with `rootCmd.AddCommand(...)`.
3. Use builders for dependencies such as `exec.NewExecutor`, `pather.NewPather`, `devcontainer.NewDevContainerCLI`, `container.NewContainerCLI`, and `vscode.NewVSCode`.
4. Normalize paths through `pather`; do not reimplement WSL path behavior.
5. Keep command output in Portuguese to match the existing CLI unless the user explicitly requests otherwise.

## Tests And Docs

1. Add tests for `{command}Impl` with existing Mockery-generated mocks where possible.
2. If the command needs a new mockable interface, update `.mockery.yml` and run `mockery --all` when available.
3. Update repository docs or command examples when user-facing behavior changes.

## Verify

Run:

```bash
go fmt ./...
go test -v ./cmd
go test -v ./...
go build -o dev main.go
```
