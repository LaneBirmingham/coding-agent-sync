# coding-agent-sync

[![CI](https://github.com/LaneBirmingham/coding-agent-sync/actions/workflows/ci.yml/badge.svg)](https://github.com/LaneBirmingham/coding-agent-sync/actions/workflows/ci.yml)
[![Latest Release](https://img.shields.io/github/v/release/LaneBirmingham/coding-agent-sync?display_name=tag)](https://github.com/LaneBirmingham/coding-agent-sync/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`coding-agent-sync` (`cas`) syncs instructions and skills across AI coding agents: Claude Code, GitHub Copilot (coding agent / agent mode), Codex, and OpenCode.

## Background

Many people switch between AI coding agents, whether that is due to subscription/usage limits or using different tools for work vs personal projects. `cas` helps keep instructions and skills in sync across Claude Code, GitHub Copilot (coding agent / agent mode), Codex, and OpenCode so you can avoid manual copy/paste between agent-specific directories and keep your best setup available everywhere.

## AI-assisted install and usage

If your coding agent supports skills, you can install the skill in this repo and let the agent handle setup and commands.

1. Add this skill to your current coding agent: [`skills/coding-agent-sync/SKILL.md`](skills/coding-agent-sync/SKILL.md).
2. Ask your agent to run setup and usage for you.
3. The skill is designed to ask for explicit consent before agent-run install steps.

Example prompt to your agent:

```text
Use the coding-agent-sync skill to install cas and sync my instructions and skills from Claude to OpenCode. Ask me before installing anything.
```

If you prefer manual control, install from releases and then ask your agent to continue:
<https://github.com/LaneBirmingham/coding-agent-sync/releases>

## Manual install and usage

### Usage commands

```bash
cas diff --from claude --to copilot,opencode --scope local
cas sync --from claude --to copilot,opencode --scope local
cas sync instructions --from claude --to opencode --scope local
cas sync skills --from claude --to copilot --scope local
cas export --from claude --scope local -o claude-local.zip
cas import --to copilot,opencode --scope local -i claude-local.zip
cas --help
cas sync --help
```
### Install (one line)

Releases: <https://github.com/LaneBirmingham/coding-agent-sync/releases>

Supported prebuilt targets:

- `darwin/arm64`
- `linux/amd64`
- `linux/arm64`

Install the latest release:

```bash
curl -fsSL https://raw.githubusercontent.com/LaneBirmingham/coding-agent-sync/main/scripts/install.sh | bash
```

Install a specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/LaneBirmingham/coding-agent-sync/main/scripts/install.sh | CAS_VERSION=v0.3.0 bash
```

Change install directory:

```bash
curl -fsSL https://raw.githubusercontent.com/LaneBirmingham/coding-agent-sync/main/scripts/install.sh | CAS_INSTALL_DIR=/usr/local/bin bash
```

Script location: [`scripts/install.sh`](scripts/install.sh)

### Install manually (source)

```bash
go install github.com/LaneBirmingham/coding-agent-sync@latest
```

### What is currently built

- Stable releases are published from `main`.
- Current prebuilt artifact targets are `darwin/arm64`, `linux/amd64`, and `linux/arm64`.
- Artifact name: `cas_<version>_<os>_<arch>.tar.gz`
- Checksum file: `SHA256SUMS`
- Other platforms can install from source with `go install`.

## Install and run locally

```bash
git clone https://github.com/LaneBirmingham/coding-agent-sync.git
cd coding-agent-sync
go run . --help
```

Build and run local binary:

```bash
go build -o cas .
./cas version
```

## Local development

```bash
go test ./...
go build ./...
go vet ./...
```

## Versioning

This project uses build-time version injection.

- `cmd/version.go` defaults to `dev` for local builds.
- Release builds inject the tagged version using `-ldflags`.

Example:

```bash
go build -ldflags "-X github.com/LaneBirmingham/coding-agent-sync/cmd.Version=1.2.3" ./...
```

## CI and releases

See [`docs/release.md`](docs/release.md) for workflow and release details.

- CI checks
- Conventional Commit enforcement
- `release-please` setup for `main` (stable)
- Binary release assets (`darwin/arm64`, `linux/amd64`, `linux/arm64`)

## Contributing

PRs are welcome. For substantial changes, please open an issue first so we can align on scope.

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for local validation steps and PR expectations.

## AI-assisted development disclaimer

This repository is developed with assistance from AI coding tools.
AI-generated suggestions may be used for code, tests, and documentation, but maintainers review, edit, and validate changes before release.

## License

MIT. See [`LICENSE`](LICENSE).
