#!/bin/sh
set -u

changed_files="$(git diff --name-only HEAD 2>/dev/null || true)"
staged_files="$(git diff --name-only --cached 2>/dev/null || true)"
untracked_files="$(git ls-files --others --exclude-standard 2>/dev/null || true)"
files="$(printf '%s\n%s\n%s\n' "$changed_files" "$staged_files" "$untracked_files" | sed '/^$/d' | sort -u)"

[ -n "$files" ] || exit 0

printf '\nDev CLI post-edit automation:\n'

printf '%s\n' '- lint: git diff --check'
if ! git diff --check; then
  printf '%s\n' '  Non-blocking: whitespace lint found issues for the agent to review.'
fi

has_go_changes=false
has_module_changes=false
has_test_changes=false

if printf '%s\n' "$files" | grep -Eq '\.go$'; then
  has_go_changes=true
fi

if printf '%s\n' "$files" | grep -Eq '(^go\.mod$|^go\.sum$)'; then
  has_module_changes=true
fi

if printf '%s\n' "$files" | grep -Eq '_test\.go$'; then
  has_test_changes=true
fi

if [ "$has_go_changes" = true ]; then
  printf '%s\n' '- lint: go fmt ./...'
  if ! go fmt ./...; then
    printf '%s\n' '  Non-blocking: go fmt failed; the agent should inspect formatting errors.'
  fi
fi

if [ "$has_go_changes" = true ] || [ "$has_module_changes" = true ]; then
  printf '%s\n' '- tidy: go mod tidy'
  if ! go mod tidy; then
    printf '%s\n' 'Dev CLI hook failure: go mod tidy failed. Inspect module/dependency changes and fix them before finishing.'
    exit 1
  fi
fi

if [ "$has_test_changes" = true ]; then
  test_packages="$(printf '%s\n' "$files" | awk '
    /_test\.go$/ {
      dir = $0
      sub("/[^/]+$", "", dir)
      if (dir == $0) {
        print "."
      } else {
        print "./" dir
      }
    }
  ' | sort -u)"

  for package in $test_packages; do
    printf '%s\n' "- tests: go test -v ${package}"
    if ! go test -v "$package"; then
      printf '%s\n' "Dev CLI hook failure: tests failed for ${package}. Inspect the failure and fix it before finishing."
      exit 1
    fi
  done
fi

printf '%s\n' 'Dev CLI post-edit automation completed.'

exit 0
