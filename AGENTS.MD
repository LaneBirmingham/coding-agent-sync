# AGENTS.MD

Guide for coding agents working in this repository.

## Project Overview

`coding-agent-sync` (`cas`) is a Go CLI that syncs instructions and skills between AI coding agents:
- Claude Code
- GitHub Copilot (coding agent / agent mode)
- Codex
- OpenCode

## Tech Stack

- Language: Go (`go 1.24.0`)
- CLI framework: Cobra (`github.com/spf13/cobra`)
- Module path: `github.com/LaneBirmingham/coding-agent-sync`

## Repository Layout

- `main.go`: CLI entrypoint.
- `cmd/`: Cobra commands (`sync`, `diff`, `export`, `import`, `version`).
- `internal/agent/`: Per-agent adapters for reading/writing instructions and skills.
- `internal/sync/`: Sync logic plus archive import/export operations.
- `internal/config/`: Agent/scope parsing and sync config types.
- `internal/archive/`: Zip archive helpers.
- `skills/coding-agent-sync/SKILL.md`: Skill instructions for AI-assisted installation and usage.
- `docs/release.md`: CI/release workflow details.

## Development Commands

Run these before considering work complete:

```bash
go test ./...
go build ./...
go vet ./...
```

Useful local runs:

```bash
go run . --help
go build -o cas .
./cas version
```

## Coding Guidelines

- Follow existing package boundaries:
  - Command parsing/flags in `cmd/*`.
  - Business logic in `internal/*`.
- Keep command behavior consistent across `sync`, `diff`, `export`, and `import`.
- Preserve existing user-facing phrasing where practical (warnings, action summaries, help text).
- Prefer table-driven tests for new logic and edge cases.
- Avoid adding new dependencies unless justified.

## Testing Expectations

- Add or update tests for behavior changes.
- Prefer focused tests near changed code:
  - `cmd/*_test.go`
  - `internal/agent/*_test.go`
  - `internal/sync/*_test.go`
  - `internal/archive/*_test.go`
- Validate both success paths and key failure/validation paths.

## CLI and Versioning Notes

- Local version defaults to `dev` (see `cmd/version.go`).
- Release builds inject version with ldflags:

```bash
go build -ldflags "-X github.com/LaneBirmingham/coding-agent-sync/cmd.Version=<version>" ./...
```

## CI and Release Notes

- CI workflow (`.github/workflows/ci.yml`) runs test/build/vet.
- PR titles to `main` must follow Conventional Commits (`.github/workflows/conventional-commits.yml`).
- Release Please on `main` manages stable releases.
- Release build publishes `darwin/arm64`, `linux/amd64`, and `linux/arm64` artifacts plus checksum file.

## Skill Maintenance Requirements

- Treat `skills/coding-agent-sync/SKILL.md` as a maintained product surface, not a one-time doc.
- If install behavior changes (`scripts/install.sh`, release artifact names, checksum format, supported platforms), update both `README.md` install instructions and `skills/coding-agent-sync/SKILL.md` in the same PR.
- If agent behavior changes (`internal/agent/*`, supported agents, scope/path precedence), update the skill's “Scope and path behavior” section in the same PR.
- Keep skill frontmatter (`name`, `description`) aligned with currently supported agents and workflows.
- When touching the skill, validate command examples still work (`cas --help`, `cas sync --help`, and `scripts/install.sh --help`).

## Commit and PR Guidance

- Use Conventional Commit style for commit messages and PR titles.
- Keep changes scoped and minimal.
- Update docs when command behavior or supported agents/scopes change.

## Agent Safety Checks

Before finishing:

1. Run `go test ./...`, `go build ./...`, and `go vet ./...`.
2. Confirm help text/examples still match implemented flags.
3. Confirm README/docs updates if user-facing behavior changed.
