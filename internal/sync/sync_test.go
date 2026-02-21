package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func localCfg(root string, from config.Agent, to []config.Agent, dryRun bool) *config.SyncConfig {
	return &config.SyncConfig{
		From:      from,
		To:        to,
		Root:      root,
		FromScope: config.ScopeLocal,
		ToScope:   config.ScopeLocal,
		DryRun:    dryRun,
	}
}

func TestSyncInstructions_ClaudeToCopilot(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# My instructions")

	cfg := localCfg(root, config.Claude, nil, false)
	action, err := SyncInstructions(cfg, config.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "synced" {
		t.Errorf("expected status synced, got %q", action.Status)
	}
	if !strings.Contains(action.String(), "synced") {
		t.Errorf("expected synced in string, got %q", action.String())
	}
	got := readFile(t, filepath.Join(root, "AGENTS.md"))
	if got != "# My instructions" {
		t.Errorf("expected '# My instructions', got %q", got)
	}
}

func TestSyncInstructions_ClaudeToCodex(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# My instructions")

	cfg := localCfg(root, config.Claude, nil, false)
	action, err := SyncInstructions(cfg, config.Codex)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "synced" {
		t.Errorf("expected status synced, got %q", action.Status)
	}
	got := readFile(t, filepath.Join(root, "AGENTS.md"))
	if got != "# My instructions" {
		t.Errorf("expected '# My instructions', got %q", got)
	}
}

func TestSyncInstructions_DryRun(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# My instructions")

	cfg := localCfg(root, config.Claude, nil, true)
	action, err := SyncInstructions(cfg, config.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "dry-run" {
		t.Errorf("expected status dry-run, got %q", action.Status)
	}
	if !strings.Contains(action.String(), "would write") {
		t.Errorf("expected dry-run message, got %q", action.String())
	}
	// File should not exist
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); !os.IsNotExist(err) {
		t.Error("expected AGENTS.md to not exist in dry-run")
	}
}

func TestSyncInstructions_SharedPath(t *testing.T) {
	root := t.TempDir()
	cfg := localCfg(root, config.Copilot, nil, false)
	action, err := SyncInstructions(cfg, config.OpenCode)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "noop" {
		t.Errorf("expected status noop, got %q", action.Status)
	}
	if !strings.Contains(action.String(), "already in sync") {
		t.Errorf("expected already-in-sync message, got %q", action.String())
	}
}

func TestSyncInstructions_SharedPath_CodexToCopilot(t *testing.T) {
	root := t.TempDir()
	cfg := localCfg(root, config.Codex, nil, false)
	action, err := SyncInstructions(cfg, config.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "noop" {
		t.Errorf("expected status noop, got %q", action.Status)
	}
	if !strings.Contains(action.String(), "already in sync") {
		t.Errorf("expected already-in-sync message, got %q", action.String())
	}
}

func TestSyncInstructions_MissingSource(t *testing.T) {
	root := t.TempDir()
	cfg := localCfg(root, config.Claude, nil, false)
	action, err := SyncInstructions(cfg, config.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "skipped" {
		t.Errorf("expected status skipped, got %q", action.Status)
	}
}

func TestSyncSkills_ClaudeToCopilot(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".claude", "skills", "my-skill", "SKILL.md"), "skill content")

	cfg := localCfg(root, config.Claude, nil, false)
	action, err := SyncSkills(cfg, config.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "synced" {
		t.Errorf("expected status synced, got %q", action.Status)
	}
	if !strings.Contains(action.String(), "synced 1 skill") {
		t.Errorf("expected synced message, got %q", action.String())
	}
	got := readFile(t, filepath.Join(root, ".github", "skills", "my-skill", "SKILL.md"))
	if got != "skill content" {
		t.Errorf("expected 'skill content', got %q", got)
	}
}

func TestSyncSkills_NoSkills(t *testing.T) {
	root := t.TempDir()
	cfg := localCfg(root, config.Claude, nil, false)
	action, err := SyncSkills(cfg, config.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "skipped" {
		t.Errorf("expected status skipped, got %q", action.Status)
	}
}

func TestSyncAll_MultipleTargets(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# Instructions")
	writeFile(t, filepath.Join(root, ".claude", "skills", "s1", "SKILL.md"), "s1")

	cfg := localCfg(root, config.Claude, []config.Agent{config.Copilot, config.OpenCode}, false)

	result, err := SyncAll(cfg, All)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Actions) != 4 {
		t.Errorf("expected 4 actions (2 targets x 2 items), got %d: %v", len(result.Actions), result.Actions)
	}
}

func TestSyncAll_InstructionsOnly(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# Instructions")

	cfg := localCfg(root, config.Claude, []config.Agent{config.Copilot}, false)

	result, err := SyncAll(cfg, Instructions)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(result.Actions))
	}
}

// --- Global sync tests ---

func TestSyncInstructions_Global_ClaudeToOpenCode(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeFile(t, filepath.Join(home, ".claude", "CLAUDE.md"), "# Global instructions")

	cfg := &config.SyncConfig{
		From:      config.Claude,
		To:        []config.Agent{config.OpenCode},
		FromScope: config.ScopeGlobal,
		ToScope:   config.ScopeGlobal,
	}

	action, err := SyncInstructions(cfg, config.OpenCode)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "synced" {
		t.Errorf("expected synced, got %q: %s", action.Status, action.Detail)
	}

	got := readFile(t, filepath.Join(home, ".config", "opencode", "AGENTS.md"))
	if got != "# Global instructions" {
		t.Errorf("expected '# Global instructions', got %q", got)
	}
}

func TestSyncInstructions_Global_ToCopilot_Skipped(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeFile(t, filepath.Join(home, ".claude", "CLAUDE.md"), "# Global")

	cfg := &config.SyncConfig{
		From:      config.Claude,
		To:        []config.Agent{config.Copilot},
		FromScope: config.ScopeGlobal,
		ToScope:   config.ScopeGlobal,
	}

	action, err := SyncInstructions(cfg, config.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "skipped" {
		t.Errorf("expected skipped, got %q: %s", action.Status, action.Detail)
	}
	if !strings.Contains(action.Detail, "does not support") {
		t.Errorf("expected unsupported message, got %q", action.Detail)
	}
}

// --- Cross-scope tests ---

func TestSyncInstructions_CrossScope_GlobalToLocal(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()

	writeFile(t, filepath.Join(home, ".claude", "CLAUDE.md"), "# Global instructions")

	cfg := &config.SyncConfig{
		From:      config.Claude,
		To:        []config.Agent{config.Claude},
		Root:      root,
		FromScope: config.ScopeGlobal,
		ToScope:   config.ScopeLocal,
	}

	action, err := SyncInstructions(cfg, config.Claude)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "synced" {
		t.Errorf("expected synced, got %q: %s", action.Status, action.Detail)
	}

	got := readFile(t, filepath.Join(root, "CLAUDE.md"))
	if got != "# Global instructions" {
		t.Errorf("expected '# Global instructions', got %q", got)
	}

	// Verify scope is shown in output
	str := action.String()
	if !strings.Contains(str, "global→local") {
		t.Errorf("expected scope info in string, got %q", str)
	}
}

func TestSyncInstructions_CrossScope_LocalToGlobal(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# Local instructions")

	cfg := &config.SyncConfig{
		From:      config.Claude,
		To:        []config.Agent{config.Claude},
		Root:      root,
		FromScope: config.ScopeLocal,
		ToScope:   config.ScopeGlobal,
	}

	action, err := SyncInstructions(cfg, config.Claude)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "synced" {
		t.Errorf("expected synced, got %q: %s", action.Status, action.Detail)
	}

	got := readFile(t, filepath.Join(home, ".claude", "CLAUDE.md"))
	if got != "# Local instructions" {
		t.Errorf("expected '# Local instructions', got %q", got)
	}
}

func TestSyncSkills_Global_ClaudeToCopilot(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeFile(t, filepath.Join(home, ".claude", "skills", "gs1", "SKILL.md"), "global skill")

	cfg := &config.SyncConfig{
		From:      config.Claude,
		To:        []config.Agent{config.Copilot},
		FromScope: config.ScopeGlobal,
		ToScope:   config.ScopeGlobal,
	}

	action, err := SyncSkills(cfg, config.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "synced" {
		t.Errorf("expected synced, got %q: %s", action.Status, action.Detail)
	}

	got := readFile(t, filepath.Join(home, ".copilot", "skills", "gs1", "SKILL.md"))
	if got != "global skill" {
		t.Errorf("expected 'global skill', got %q", got)
	}
}

func TestSyncInstructions_Global_ClaudeToCodex(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeFile(t, filepath.Join(home, ".claude", "CLAUDE.md"), "# Global instructions")

	cfg := &config.SyncConfig{
		From:      config.Claude,
		To:        []config.Agent{config.Codex},
		FromScope: config.ScopeGlobal,
		ToScope:   config.ScopeGlobal,
	}

	action, err := SyncInstructions(cfg, config.Codex)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "synced" {
		t.Errorf("expected synced, got %q: %s", action.Status, action.Detail)
	}

	got := readFile(t, filepath.Join(home, ".codex", "AGENTS.md"))
	if got != "# Global instructions" {
		t.Errorf("expected '# Global instructions', got %q", got)
	}
}

func TestSyncSkills_Global_ClaudeToCodex(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeFile(t, filepath.Join(home, ".claude", "skills", "gs1", "SKILL.md"), "global skill")

	cfg := &config.SyncConfig{
		From:      config.Claude,
		To:        []config.Agent{config.Codex},
		FromScope: config.ScopeGlobal,
		ToScope:   config.ScopeGlobal,
	}

	action, err := SyncSkills(cfg, config.Codex)
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "synced" {
		t.Errorf("expected synced, got %q: %s", action.Status, action.Detail)
	}

	got := readFile(t, filepath.Join(home, ".agents", "skills", "gs1", "SKILL.md"))
	if got != "global skill" {
		t.Errorf("expected 'global skill', got %q", got)
	}
}

func TestSyncAction_String_ShowsScope(t *testing.T) {
	action := SyncAction{
		Kind:      Instructions,
		From:      config.Claude,
		To:        config.Copilot,
		FromScope: config.ScopeGlobal,
		ToScope:   config.ScopeLocal,
		Status:    "synced",
		Detail:    "synced (42 bytes)",
	}
	got := action.String()
	if !strings.Contains(got, "global→local") {
		t.Errorf("expected scope in output, got %q", got)
	}
}

func TestSyncAction_String_HidesScopeForLocal(t *testing.T) {
	action := SyncAction{
		Kind:      Instructions,
		From:      config.Claude,
		To:        config.Copilot,
		FromScope: config.ScopeLocal,
		ToScope:   config.ScopeLocal,
		Status:    "synced",
		Detail:    "synced (42 bytes)",
	}
	got := action.String()
	if strings.Contains(got, "local→local") {
		t.Errorf("should not show scope for local→local, got %q", got)
	}
}
