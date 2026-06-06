# Dev CLI Codex Harness

This directory contains repo-local Codex skills for Dev CLI maintenance. Harness instructions are written in English so agents can follow one consistent workflow, while existing CLI output, tests, and user-facing Portuguese copy remain unchanged unless a task explicitly changes them.

## Available Skills

- `release-version`: prepare and publish a version release, including release notes, verification, tag push, and GitHub Release creation or update.
- `new-internal-package`: scaffold a new package under `internal/` using the repository interface, builder, implementation, test, and mock pattern.
- `fix-bug`: reproduce a defect, patch it narrowly, add regression coverage, and run targeted plus full Go verification.
- `add-command`: add a Cobra command under `cmd/` with dependency injection, registration, tests, and documentation updates when behavior changes.
- `commit-changes`: split worktree changes into focused commits with strict Conventional Commit messages.

## Trigger Examples

Use explicit invocation when you want the workflow to drive the task:

```text
$fix-bug dev open fails for WSL paths
$add-command add a status command
$new-internal-package create internal/update
$release-version prepare v1.4.0
$commit-changes commit the current changes in logical parts
```

Codex may also invoke a skill implicitly when a request clearly matches the skill description.

## Verification

The normal repository checks are:

```bash
go fmt ./...
go test -v ./...
go build -o dev main.go
```

Release tasks also inspect tags/releases with `git` and `gh`, and should only push tags or mutate GitHub Releases after approval.
