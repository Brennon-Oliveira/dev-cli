---
name: commit-changes
description: "Use when Codex needs to inspect the current git worktree, split changes into focused commits, stage files or hunks selectively, and create commit messages that strictly follow Conventional Commits 1.0.0."
---

# Commit Changes

Use this workflow to commit repository changes cleanly and in small logical groups.

## Inspect

1. Run `git status --short`.
2. Inspect tracked diffs with `git diff` and staged diffs with `git diff --cached`.
3. Inspect untracked files before staging them. Use `find`, `rg --files`, `sed`, or similar read-only commands.
4. Identify independent themes. Prefer multiple commits when changes serve different purposes.
5. Do not stage unrelated user changes unless the user explicitly asks to include them.

## Stage

1. Stage only one logical theme at a time.
2. Use explicit paths with `git add <path>...`.
3. When one file contains changes for multiple themes, use a patch staging command or ask before mixing unrelated changes.
4. Verify the staged content with `git diff --cached --stat` and `git diff --cached`.
5. Commit only after the staged diff matches the intended theme.

## Message Standard

Use this self-contained Conventional Commits subset. Do not look up external documentation for routine commits.

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Rules:

- Prefix every commit with `type`, optional `(scope)`, optional `!`, then `: `.
- Use `feat` for a new feature or capability.
- Use `fix` for a bug fix.
- Use `ci` for CI/workflow changes, `docs` for documentation, `test` for tests, `refactor` for behavior-preserving code changes, `style` for formatting-only changes, `perf` for performance improvements, `build` for build/dependency changes, and `chore` for maintenance that fits none of the above.
- Use a short noun scope when it clarifies the affected area, for example `harness`, `release`, `cmd`, `config`, or `devcontainer`.
- Write the description as a short imperative summary, lowercase unless a proper noun requires casing, with no trailing punctuation.
- Mark breaking changes with `!` before `:` or a footer beginning exactly `BREAKING CHANGE: `.
- Add a blank line before a body and before footers.
- Add a body only when the staged diff needs non-obvious context.

## Commit

1. Choose the narrowest type and scope that match the staged diff.
2. Run `git commit -m "<type>(<scope>): <description>"` for simple commits.
3. For commits needing a body or footer, use multiple `-m` arguments.
4. After each commit, run `git status --short` and continue with the next logical theme.
5. At the end, report each commit hash and message, plus any remaining uncommitted files.
