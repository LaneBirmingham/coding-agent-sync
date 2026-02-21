package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/LaneBirmingham/coding-agent-sync/internal/agent"
	"github.com/LaneBirmingham/coding-agent-sync/internal/archive"
	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

func TestExportDryRun(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(t.TempDir(), "export.zip")

	if err := os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte("# Instructions"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".claude", "skills", "s1"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".claude", "skills", "s1", "SKILL.md"), []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := Export(&config.ExportConfig{
		From:       config.Claude,
		Root:       root,
		Scope:      config.ScopeLocal,
		Output:     output,
		DryRun:     true,
		CASVersion: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(result.Actions))
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("expected no archive file in dry-run, got err=%v", err)
	}
}

func TestExportAndImportRoundTrip(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(t.TempDir(), "roundtrip.zip")

	if err := os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte("# A"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".claude", "skills", "skill-a"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".claude", "skills", "skill-a", "SKILL.md"), []byte("skill-a"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Export(&config.ExportConfig{
		From:       config.Claude,
		Root:       root,
		Scope:      config.ScopeLocal,
		Output:     archivePath,
		DryRun:     false,
		CASVersion: "test",
	}); err != nil {
		t.Fatal(err)
	}

	result, err := Import(&config.ImportConfig{
		To:     []config.Agent{config.Copilot},
		Root:   root,
		Scope:  config.ScopeLocal,
		Input:  archivePath,
		DryRun: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(result.Actions))
	}
	if len(result.Warnings) == 0 {
		t.Fatal("expected mismatch warning when importing claude archive to copilot")
	}

	inst, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(inst) != "# A" {
		t.Fatalf("expected imported instructions '# A', got %q", string(inst))
	}
	skill, err := os.ReadFile(filepath.Join(root, ".github", "skills", "skill-a", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(skill) != "skill-a" {
		t.Fatalf("expected imported skill content, got %q", string(skill))
	}
}

func TestImportSkipsUnsupportedInstructions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	archivePath := filepath.Join(t.TempDir(), "global.zip")
	if err := archive.Write(archivePath, &archive.Archive{
		Manifest: &archive.Manifest{
			Version:    archive.FormatVersion,
			Agent:      "claude",
			Scope:      "global",
			ExportedAt: time.Now().UTC(),
		},
		Instructions: &agent.Instruction{Content: "# Global"},
		Skills:       []agent.Skill{{Name: "shared", Content: "shared"}},
	}); err != nil {
		t.Fatal(err)
	}

	result, err := Import(&config.ImportConfig{
		To:     []config.Agent{config.Copilot},
		Root:   t.TempDir(),
		Scope:  config.ScopeGlobal,
		Input:  archivePath,
		DryRun: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(result.Actions))
	}
	if result.Actions[0].Status != "skipped" {
		t.Fatalf("expected instructions action to be skipped, got %q", result.Actions[0].Status)
	}
	if !strings.Contains(result.Actions[0].Detail, "does not support") {
		t.Fatalf("expected unsupported detail, got %q", result.Actions[0].Detail)
	}
	if result.Actions[1].Status != "imported" {
		t.Fatalf("expected skills action imported, got %q", result.Actions[1].Status)
	}

	skillPath := filepath.Join(home, ".copilot", "skills", "shared", "SKILL.md")
	data, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "shared" {
		t.Fatalf("expected imported skill content, got %q", string(data))
	}
}
