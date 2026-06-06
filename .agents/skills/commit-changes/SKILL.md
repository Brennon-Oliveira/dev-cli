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

Follow Conventional Commits 1.0.0 exactly:

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Rules:

- Prefix every commit with a type, optional scope, optional `!`, then `: `.
- Use `feat` only when adding a user-facing feature or capability.
- Use `fix` only when fixing a bug.
- Use allowed supporting types when appropriate, such as `docs`, `ci`, `chore`, `test`, `refactor`, `style`, `perf`, or `build`.
- Use a noun scope when it helps identify the codebase area, for example `chore(harness): ...`.
- Keep the description short, imperative, lowercase unless a proper noun requires casing, and do not end it with punctuation.
- For breaking changes, add `!` before the colon or include a footer beginning exactly with `BREAKING CHANGE: `.
- Put a blank line before a body and before footers.
- Use body text only when the why/risk is important and not obvious from the diff.

## Commit

1. Choose the narrowest type and scope that match the staged diff.
2. Run `git commit -m "<type>(<scope>): <description>"` for simple commits.
3. For commits needing a body or footer, use multiple `-m` arguments.
4. After each commit, run `git status --short` and continue with the next logical theme.
5. At the end, report each commit hash and message, plus any remaining uncommitted files.

Source for the message format: https://www.conventionalcommits.org/en/v1.0.0/
