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

func TestSyncInstructions_ClaudeToCopilot(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# My instructions")

	action, err := SyncInstructions(root, config.Claude, config.Copilot, false)
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

func TestSyncInstructions_DryRun(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# My instructions")

	action, err := SyncInstructions(root, config.Claude, config.Copilot, true)
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
	action, err := SyncInstructions(root, config.Copilot, config.OpenCode, false)
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
	action, err := SyncInstructions(root, config.Claude, config.Copilot, false)
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

	action, err := SyncSkills(root, config.Claude, config.Copilot, false)
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
	action, err := SyncSkills(root, config.Claude, config.Copilot, false)
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

	cfg := &config.SyncConfig{
		From: config.Claude,
		To:   []config.Agent{config.Copilot, config.OpenCode},
		Root: root,
	}

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

	cfg := &config.SyncConfig{
		From: config.Claude,
		To:   []config.Agent{config.Copilot},
		Root: root,
	}

	result, err := SyncAll(cfg, Instructions)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(result.Actions))
	}
}
