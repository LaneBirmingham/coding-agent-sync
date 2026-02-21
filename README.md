# coding-agent-sync

`coding-agent-sync` (`cas`) syncs instructions and skills between Claude Code, GitHub Copilot, and OpenCode.

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
- `release-please` setup for `main` (stable)
- Binary release assets (`darwin/arm64`)
