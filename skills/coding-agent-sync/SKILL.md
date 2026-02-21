---
name: coding-agent-sync
description: Install and operate the `cas` (coding-agent-sync) CLI to sync instructions and skills between Claude Code, GitHub Copilot, Codex, and OpenCode. Use when asked to migrate, compare, back up, or standardize agent instructions/skills across agents or scopes, including downloading/installing the binary when `cas` is missing.
---

# coding-agent-sync

Install `cas` first, then run a dry-run preview before writing any changes.

## 1) Ensure `cas` is installed

Check whether `cas` exists:

```bash
cas version
```

If missing, request explicit consent before any install action:

Ask:
"`cas` is not installed. I can install it for you automatically (to `~/.local/bin`) from the official releases page: https://github.com/LaneBirmingham/coding-agent-sync/releases . If you prefer safer manual control, you can install it yourself from that same link and then ask me to continue. Do you want me to install it now?"

If the user declines or prefers manual install, share the releases link and stop until they confirm installation.
If the user approves, install in this order.

### Preferred: official install script (release binaries)

```bash
curl -fsSL https://raw.githubusercontent.com/LaneBirmingham/coding-agent-sync/main/scripts/install.sh | bash
```

Optional script overrides:

```bash
curl -fsSL https://raw.githubusercontent.com/LaneBirmingham/coding-agent-sync/main/scripts/install.sh | CAS_VERSION=v0.3.0 bash
curl -fsSL https://raw.githubusercontent.com/LaneBirmingham/coding-agent-sync/main/scripts/install.sh | CAS_INSTALL_DIR=/usr/local/bin bash
```

Release binaries currently cover `darwin/arm64`, `linux/amd64`, and `linux/arm64`.

### Fallback: build from source (`go install`)

```bash
go install github.com/LaneBirmingham/coding-agent-sync@latest
cas version
```

If `cas` is still not found after `go install`, add `$(go env GOPATH)/bin` to `PATH`.
If auto-install fails because of permissions or environment restrictions, provide manual install guidance and link:
`https://github.com/LaneBirmingham/coding-agent-sync/releases`

## 2) Use safe workflow by default

Run dry-run first:

```bash
cas diff --from <agent> --to <agent[,agent...]> --scope <local|global>
```

Then apply:

```bash
cas sync --from <agent> --to <agent[,agent...]> --scope <local|global>
```

For partial sync:

```bash
cas sync instructions --from claude --to opencode --scope local
cas sync skills --from claude --to copilot --scope local
cas sync --from codex --to claude,opencode --scope local
```

## 3) Use archive workflow for migration/backup

Export:

```bash
cas export --from claude --scope local -o claude-local.zip
```

Import:

```bash
cas import --to copilot,opencode --scope local -i claude-local.zip
```

Preview archive operations with `--dry-run` before writing.

## 4) Scope and path behavior

Use `local` for project files, `global` for user-level config.

Local targets:

- Claude instructions: `CLAUDE.md` (or `.claude/CLAUDE.md` as source)
- Copilot instructions: `AGENTS.md`
- Codex instructions (read precedence): `AGENTS.override.md`, `AGENTS.md`, `TEAM_GUIDE.md`, `.agents.md`
- OpenCode instructions: `AGENTS.md`
- Claude skills: `.claude/skills/*/SKILL.md`
- Copilot skills: `.github/skills/*/SKILL.md`
- Codex skills (read): `.agents/skills/*/SKILL.md` (fallback `.codex/skills/*/SKILL.md`)
- OpenCode skills: `.opencode/skills/*/SKILL.md`

Global targets:

- Claude instructions: `~/.claude/CLAUDE.md`
- Copilot skills: `~/.copilot/skills/*/SKILL.md`
- Codex instructions (read): `$CODEX_HOME/AGENTS.override.md`, `$CODEX_HOME/AGENTS.md` (fallback `~/.codex/*`); write target is `$CODEX_HOME/AGENTS.md`
- Codex skills (read): `~/.agents/skills/*/SKILL.md` (fallback `$CODEX_HOME/skills/*/SKILL.md`); write target is `~/.agents/skills/*/SKILL.md`
- OpenCode instructions: `~/.config/opencode/AGENTS.md`
- OpenCode skills: `~/.config/opencode/skills/*/SKILL.md`

Do not attempt Copilot global instructions; Copilot global instructions are unsupported.
