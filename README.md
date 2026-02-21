# coding-agent-sync

`coding-agent-sync` (`cas`) syncs instructions and skills between Claude Code, GitHub Copilot, and OpenCode.

## Install

### Download prebuilt binary (currently `darwin/arm64`)

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

### Build/install from source (all Go-supported platforms)

```bash
go install github.com/LaneBirmingham/coding-agent-sync@latest
```

## CLI Quick Start

Preview changes before writing:

```bash
cas diff --from claude --to copilot,opencode --scope local
```

Sync only skills:

```bash
cas sync skills --from claude --to copilot --from-scope local --to-scope local
```

Export and import:

```bash
cas export --from claude --scope local -o claude-local.zip
cas import --to copilot,opencode --scope local -i claude-local.zip
```

## Local Development

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

## CI and Releases

See [`docs/release.md`](docs/release.md) for:

- CI checks
- Conventional Commit enforcement
- `release-please` setup for `main` (stable) and `dev` (beta prerelease)
- Binary release assets (`darwin/arm64`)

## Included Skill

This repo includes a Codex skill at [`skills/coding-agent-sync/SKILL.md`](skills/coding-agent-sync/SKILL.md) that teaches an agent to:

- install `cas` (binary-first, source fallback)
- run safe `diff`/`sync` workflows
- use `export`/`import` for migration and backup
