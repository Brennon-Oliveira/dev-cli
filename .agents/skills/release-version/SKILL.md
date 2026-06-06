---
name: release-version
description: Use when preparing or publishing a Dev CLI version release, including inspecting git tags and GitHub Releases, drafting English release notes, running verification, creating and pushing a version tag, and creating or updating the GitHub Release.
---

# Release Version

Use this workflow for Dev CLI releases. Keep release notes and harness output in English. Do not translate existing Portuguese CLI messages.

## Preconditions

1. Confirm the intended version tag, for example `v1.4.0`.
2. Check the worktree with `git status --short`; do not release from an unintentionally dirty tree.
3. Inspect existing tags and releases:
   - `git describe --tags --abbrev=0` for the latest local tag.
   - `git log --oneline <latest-tag>..HEAD` for unreleased commits.
   - `gh release list` and `gh release view <tag>` when GitHub state matters.
4. Remember that `.github/workflows/build.yml` runs on `v*` tag pushes and creates or updates the GitHub Release. `.github/workflows/snap-publish.yml` runs on the `release: released` event, so publish only after the tag and release notes are correct.

## Release Notes

Draft standardized release notes before mutating remote state:

```md
## What's Changed

- Concise user-facing or maintainer-facing change.

## Verification

- `go test -v ./...`
- `go build -o dev main.go`
```

If there are no meaningful changes since the previous tag, stop and ask whether to continue.

## Verification

Run:

```bash
go test -v ./...
go build -o dev main.go
```

Run `go fmt ./...` first if Go files changed and formatting has not already been checked.

## Publish

1. Create the local tag only after verification:
   ```bash
   git tag vX.Y.Z
   ```
2. Push the tag only with approval:
   ```bash
   git push origin vX.Y.Z
   ```
3. Create or update the GitHub Release only with approval:
   ```bash
   gh release create vX.Y.Z --title "vX.Y.Z" --notes-file <notes-file>
   gh release edit vX.Y.Z --title "vX.Y.Z" --notes-file <notes-file>
   ```
4. After the workflow finishes, inspect release assets with `gh release view vX.Y.Z`.

Do not use broad shell prefixes or destructive git commands for release work.
