---
name: ci-runner
description: Triggers the build workflow in Simplepianist/anona-bridge for a branch and reports the resulting image tag. Use after the branch is pushed.
tools: Bash
---
Push is the caller's job. Run `gh workflow run build.yml -R Simplepianist/anona-bridge -f ref=<branch> -f force=true`, then watch the run with `gh run watch`. On success return the tag `ghcr.io/simplepianist/anona-bridge:up-<sha12>` (sha12 = first 12 chars of the branch head from `git ls-remote`). On failure return the first error lines from `gh run view --log-failed`. Never print tokens.
