# coding-agent-sync

`coding-agent-sync` (`cas`) syncs instructions and skills between Claude Code, GitHub Copilot, and OpenCode.

## Background

Many teams keep the same guidance in multiple places (`CLAUDE.md`, `AGENTS.md`, and skill folders). `cas` helps you keep these files aligned with a dry-run-first workflow so you can preview changes before writing.

## Install

### Option 1: Download a release (recommended)

Releases: <https://github.com/LaneBirmingham/coding-agent-sync/releases>

### What is currently built

- Stable releases are published from `main`.
- `dev` runs CI but does not publish beta releases.
- Current prebuilt artifact target is `darwin/arm64`.
- Artifact name: `cas_<version>_darwin_arm64.tar.gz`
- Checksum file: `SHA256SUMS`
- Other platforms can install from source with `go install`.

Quick install for macOS Apple Silicon:

```bash
VERSION="$(curl -fsSL https://api.github.com/repos/LaneBirmingham/coding-agent-sync/releases/latest | sed -n 's/.*"tag_name": "v\([^"]*\)".*/\1/p')"
ASSET="cas_${VERSION}_darwin_arm64.tar.gz"
curl -fL "https://github.com/LaneBirmingham/coding-agent-sync/releases/download/v${VERSION}/${ASSET}" -o "/tmp/${ASSET}"
tar -xzf "/tmp/${ASSET}" -C /tmp
mkdir -p "${HOME}/.local/bin"
install -m 0755 "/tmp/cas_${VERSION}_darwin_arm64" "${HOME}/.local/bin/cas"
export PATH="${HOME}/.local/bin:${PATH}"
cas version
```

### Option 2: Install from source

```bash
go install github.com/LaneBirmingham/coding-agent-sync@latest
```

## Usage

Preview before applying changes:

```bash
cas diff --from claude --to copilot,opencode --scope local
cas sync --from claude --to copilot,opencode --scope local
```

Sync only one content type:

```bash
cas sync instructions --from claude --to opencode --scope local
cas sync skills --from claude --to copilot --scope local
```

Use archive export/import for migration or backup:

```bash
cas export --from claude --scope local -o claude-local.zip
cas import --to copilot,opencode --scope local -i claude-local.zip
```

Show command help:

```bash
cas --help
cas sync --help
```

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

## Included skill

This repo includes a Codex skill at [`skills/coding-agent-sync/SKILL.md`](skills/coding-agent-sync/SKILL.md) that teaches an agent to install `cas` and run safe sync/export/import workflows.

## AI-assisted development disclaimer

This repository is developed with assistance from AI coding tools.
AI-generated suggestions may be used for code, tests, and documentation, but maintainers review, edit, and validate changes before release.

## License

MIT. See [`LICENSE`](LICENSE).
