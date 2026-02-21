# Build and Release Process

## Overview

This repository uses GitHub Actions + release-please.

- `main` drives stable releases.
- Conventional Commits drive changelog and semantic version bumps.

## Workflows

- `.github/workflows/ci.yml`
  - Runs on PRs to `main`, and pushes to `main`/`dev`.
  - Executes:
    - `go test ./...`
    - `go build ./...`
    - `go vet ./...`

- `.github/workflows/conventional-commits.yml`
  - Enforces Conventional Commit style PR titles.
  - Use squash merges so the PR title becomes the release-relevant commit message.

- `.github/workflows/release-please-main.yml`
  - Runs release-please against `main`.
  - Creates/updates release PRs and publishes stable releases.

- `.github/workflows/release-build.yml`
  - Triggered on GitHub release publish.
  - Builds binaries for:
    - `darwin/arm64`
    - `linux/amd64`
    - `linux/arm64`
  - Injects version at build time with ldflags.
  - Uploads `.tar.gz` assets and `SHA256SUMS`.

## Release-Please Configuration

- Stable config:
  - `.release-please-config-main.json`
  - `.release-please-manifest-main.json`

## Required Repository Setup

1. Add repository secret:
   - `RELEASE_PLEASE_TOKEN` (PAT with repo permissions).

2. In repository settings, allow Actions to create and approve pull requests if required by your org policy.

3. Protect `main` branch and require checks:
   - `ci / test-build-vet`
   - `conventional-commits / semantic-pr-title`

4. For `dev`, if you continue pushing directly (no PRs), keep CI on push and skip PR-only protections.

## Commit and PR conventions

Use Conventional Commit PR titles, for example:

- `feat: add opencode export command`
- `fix: handle missing skills directory`
- `docs: clarify sync scope behavior`

## Optional local pre-commit hook

A basic hook is available at `.githooks/pre-commit` and runs `go vet ./...`.

Enable it locally:

```bash
git config core.hooksPath .githooks
```
