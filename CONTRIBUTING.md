# Contributing

Thanks for contributing to `coding-agent-sync`.

## Before opening a PR

1. Open an issue first for substantial changes (new commands, behavior changes, or release/process changes).
2. Keep changes focused and avoid unrelated refactors.
3. Add/update tests for behavior changes.
4. Run local validation:

```bash
go test ./...
go build ./...
go vet ./...
```

## PR requirements

- Use a Conventional Commit style PR title, for example:
  - `feat: add linux release artifact`
  - `fix: handle missing global instructions path`
  - `docs: clarify install instructions`
- Update docs (`README.md`, `docs/release.md`) when user-facing behavior changes.
- Keep CI green.

## Release notes

Releases are managed by release-please from `main`. See `docs/release.md` for details.
