---
name: fix-bug
description: "Use when fixing a Dev CLI bug: reproduce the issue, isolate the failing command or internal package, patch narrowly, add regression coverage, and run targeted plus full Go verification."
---

# Fix Bug

Use this workflow for defect fixes.

## Reproduce And Isolate

1. Capture the reported behavior and expected behavior in concrete terms.
2. Find the smallest affected surface:
   - `cmd/` for Cobra command wiring or argument behavior.
   - `internal/devcontainer` and `internal/pather` for path and devcontainer resolution.
   - `internal/container`, `internal/exec`, or `internal/vscode` for external command behavior.
3. Reproduce with an existing test when possible. If no test exists, add a failing regression test before patching when practical.

## Patch

1. Keep the change narrow and local to the failing behavior.
2. Preserve Portuguese CLI/user-facing output unless the bug is specifically copy-related.
3. Use existing builders, interfaces, and mocks.
4. Avoid unrelated refactors or generated file churn.

## Verify

Run the targeted package or command tests first, for example:

```bash
go test -v ./internal/pather
go test -v ./cmd
```

Then run:

```bash
go test -v ./...
go build -o dev main.go
```

If the bug touched path translation, devcontainer parsing, config persistence, or command execution, include that risk area in the final summary.
