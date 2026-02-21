# CI and Release Guide

This document describes the CI checks, release automation, and binary publishing flow for `coding-agent-sync`.

## Workflow Overview

| Workflow | File | Trigger | Purpose |
|---|---|---|---|
| CI | `.github/workflows/ci.yml` | Pull requests, push to `main`/`dev` | Run `go test`, `go build`, `go vet` |
| Conventional Commits | `.github/workflows/conventional-commits.yml` | `pull_request_target` events | Enforce semantic PR title |
| Release Please (stable) | `.github/workflows/release-please-main.yml` | Push to `main` | Manage stable release PR + tags |
| Release Please (beta) | `.github/workflows/release-please-dev.yml` | Push to `dev` | Manage prerelease PR + tags |
| Release Build | `.github/workflows/release-build.yml` | GitHub release `published` or `prereleased` | Build and upload release artifacts |

## CI Checks

`ci.yml` runs:

```bash
go test ./...
go build ./...
go vet ./...
```

Use this locally before opening a PR:

```bash
go test ./... && go build ./... && go vet ./...
```

## Conventional Commit Enforcement

`conventional-commits.yml` checks PR titles with `amannn/action-semantic-pull-request@v5`.

Accepted types:

- `feat`
- `fix`
- `docs`
- `style`
- `refactor`
- `perf`
- `test`
- `build`
- `ci`
- `chore`
- `revert`

## Release Automation

Two `release-please` configs are used:

- `main` branch: stable tags like `v1.2.3`
- `dev` branch: prerelease tags like `v1.2.4-beta.1`

Config files:

- `.release-please-config-main.json`
- `.release-please-config-dev.json`
- `.release-please-manifest-main.json`
- `.release-please-manifest-dev.json`

### Required secret

- `RELEASE_PLEASE_TOKEN` with permission to create and update PRs/issues/releases.

### Release flow

1. Merge Conventional Commit PRs into `main` (stable) or `dev` (beta).
2. `release-please` updates or opens a release PR.
3. Merge the release PR.
4. `release-please` creates a GitHub release and tag.
5. `release-build.yml` runs on that release and uploads binary artifacts.

## Published Binary Artifacts

Current release build publishes:

- `cas_<version>_darwin_arm64.tar.gz`
- `SHA256SUMS`

Build flags include version injection:

```bash
-ldflags "-s -w -X github.com/LaneBirmingham/coding-agent-sync/cmd.Version=${VERSION}"
```

### Verify and install

```bash
VERSION="0.1.0"
ASSET="cas_${VERSION}_darwin_arm64.tar.gz"
TMPDIR="$(mktemp -d)"
curl -fL "https://github.com/LaneBirmingham/coding-agent-sync/releases/download/v${VERSION}/${ASSET}" -o "${ASSET}"
curl -fL "https://github.com/LaneBirmingham/coding-agent-sync/releases/download/v${VERSION}/SHA256SUMS" -o SHA256SUMS
EXPECTED="$(awk "/${ASSET}/{print \$1; exit}" SHA256SUMS)"
ACTUAL="$(shasum -a 256 "${ASSET}" | awk '{print $1}')"
test "${EXPECTED}" = "${ACTUAL}"
tar -xzf "${ASSET}" -C "${TMPDIR}"
mkdir -p "${HOME}/.local/bin"
install -m 0755 "${TMPDIR}/cas_${VERSION}_darwin_arm64" "${HOME}/.local/bin/cas"
export PATH="${HOME}/.local/bin:${PATH}"
cas version
```

## Running `cas` in CI

Example workflow snippet:

```yaml
- name: Install cas
  run: |
    go install github.com/LaneBirmingham/coding-agent-sync@latest
    echo "$(go env GOPATH)/bin" >> "$GITHUB_PATH"

- name: Preview sync
  run: cas diff --from claude --to copilot,opencode --scope local

- name: Apply sync
  run: cas sync --from claude --to copilot,opencode --scope local
```
