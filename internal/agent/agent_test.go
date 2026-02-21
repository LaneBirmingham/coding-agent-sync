package agent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// --- Claude local tests ---

func TestClaude_ReadInstructions_RootFile(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, "CLAUDE.md"), "# Instructions")

	c := &Claude{}
	inst, err := c.ReadInstructions(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if inst == nil || inst.Content != "# Instructions" {
		t.Errorf("expected '# Instructions', got %v", inst)
	}
}

func TestClaude_ReadInstructions_DotClaudeFile(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, ".claude", "CLAUDE.md"), "# From .claude")
	writeTestFile(t, filepath.Join(root, "CLAUDE.md"), "# From root")

	c := &Claude{}
	inst, err := c.ReadInstructions(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if inst.Content != "# From .claude" {
		t.Errorf("expected .claude/CLAUDE.md to take priority, got %q", inst.Content)
	}
}

func TestClaude_ReadInstructions_Missing(t *testing.T) {
	root := setupTestDir(t)
	c := &Claude{}
	inst, err := c.ReadInstructions(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if inst != nil {
		t.Error("expected nil for missing instructions")
	}
}

func TestClaude_WriteInstructions(t *testing.T) {
	root := setupTestDir(t)
	c := &Claude{}
	err := c.WriteInstructions(config.Local(root), &Instruction{Content: "# Test"})
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(root, "CLAUDE.md"))
	if got != "# Test" {
		t.Errorf("expected '# Test', got %q", got)
	}
}

func TestClaude_InstructionsPath(t *testing.T) {
	c := &Claude{}
	got := c.InstructionsPath(config.Local("/project"))
	want := filepath.Join("/project", "CLAUDE.md")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestClaude_ReadWriteSkills(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, ".claude", "skills", "my-skill", "SKILL.md"), "skill content")

	c := &Claude{}
	skills, err := c.ReadSkills(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "my-skill" || skills[0].Content != "skill content" {
		t.Errorf("unexpected skills: %v", skills)
	}

	root2 := setupTestDir(t)
	err = c.WriteSkills(config.Local(root2), skills)
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(root2, ".claude", "skills", "my-skill", "SKILL.md"))
	if got != "skill content" {
		t.Errorf("expected 'skill content', got %q", got)
	}
}

// --- Claude global tests ---

func TestClaude_Global_ReadWriteInstructions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	c := &Claude{}

	// No file yet — should return nil
	inst, err := c.ReadInstructions(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if inst != nil {
		t.Error("expected nil for missing global instructions")
	}

	// Write global instructions
	err = c.WriteInstructions(config.Global(), &Instruction{Content: "# Global Claude"})
	if err != nil {
		t.Fatal(err)
	}

	// Read back
	inst, err = c.ReadInstructions(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if inst == nil || inst.Content != "# Global Claude" {
		t.Errorf("expected '# Global Claude', got %v", inst)
	}

	// Verify path
	path := c.InstructionsPath(config.Global())
	if path != filepath.Join(home, ".claude", "CLAUDE.md") {
		t.Errorf("unexpected global path: %q", path)
	}
}

func TestClaude_Global_ReadWriteSkills(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	c := &Claude{}
	writeTestFile(t, filepath.Join(home, ".claude", "skills", "global-skill", "SKILL.md"), "global skill")

	skills, err := c.ReadSkills(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "global-skill" {
		t.Errorf("unexpected skills: %v", skills)
	}
}

// --- Copilot local tests ---

func TestCopilot_ReadWriteInstructions(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, "AGENTS.md"), "# Copilot instructions")

	c := &Copilot{}
	inst, err := c.ReadInstructions(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if inst.Content != "# Copilot instructions" {
		t.Errorf("unexpected: %q", inst.Content)
	}

	root2 := setupTestDir(t)
	err = c.WriteInstructions(config.Local(root2), inst)
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(root2, "AGENTS.md"))
	if got != "# Copilot instructions" {
		t.Errorf("expected '# Copilot instructions', got %q", got)
	}
}

func TestCopilot_Skills(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, ".github", "skills", "s1", "SKILL.md"), "s1 content")

	c := &Copilot{}
	skills, err := c.ReadSkills(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "s1" {
		t.Errorf("unexpected skills: %v", skills)
	}
}

// --- Copilot global tests ---

func TestCopilot_Global_InstructionsUnsupported(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	c := &Copilot{}

	// InstructionsPath should return "" for global
	path := c.InstructionsPath(config.Global())
	if path != "" {
		t.Errorf("expected empty path for unsupported global instructions, got %q", path)
	}

	// ReadInstructions should return nil, nil
	inst, err := c.ReadInstructions(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if inst != nil {
		t.Error("expected nil for unsupported global instructions")
	}

	// WriteInstructions should return error
	err = c.WriteInstructions(config.Global(), &Instruction{Content: "test"})
	if err == nil {
		t.Error("expected error for unsupported global instructions write")
	}
}

func TestCopilot_Global_Skills(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeTestFile(t, filepath.Join(home, ".copilot", "skills", "gs1", "SKILL.md"), "global copilot skill")

	c := &Copilot{}
	skills, err := c.ReadSkills(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "gs1" {
		t.Errorf("unexpected skills: %v", skills)
	}
}

// --- OpenCode local tests ---

func TestOpenCode_ReadWriteInstructions(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, "AGENTS.md"), "# OC instructions")

	o := &OpenCode{}
	inst, err := o.ReadInstructions(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if inst.Content != "# OC instructions" {
		t.Errorf("unexpected: %q", inst.Content)
	}
}

func TestOpenCode_Skills(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, ".opencode", "skills", "s1", "SKILL.md"), "s1 content")

	o := &OpenCode{}
	skills, err := o.ReadSkills(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "s1" {
		t.Errorf("unexpected skills: %v", skills)
	}
}

// --- OpenCode global tests ---

func TestOpenCode_Global_ReadInstructions_Canonical(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeTestFile(t, filepath.Join(home, ".config", "opencode", "AGENTS.md"), "# OC global")

	o := &OpenCode{}
	inst, err := o.ReadInstructions(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if inst == nil || inst.Content != "# OC global" {
		t.Errorf("expected '# OC global', got %v", inst)
	}
}

func TestOpenCode_Global_ReadInstructions_Fallback(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// No canonical file, but Claude's global config exists
	writeTestFile(t, filepath.Join(home, ".claude", "CLAUDE.md"), "# Claude global fallback")

	o := &OpenCode{}
	inst, err := o.ReadInstructions(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if inst == nil || inst.Content != "# Claude global fallback" {
		t.Errorf("expected fallback content, got %v", inst)
	}
}

func TestOpenCode_Global_WriteInstructions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	o := &OpenCode{}
	err := o.WriteInstructions(config.Global(), &Instruction{Content: "# Written"})
	if err != nil {
		t.Fatal(err)
	}

	// Should write to canonical path, not Claude's path
	got := readTestFile(t, filepath.Join(home, ".config", "opencode", "AGENTS.md"))
	if got != "# Written" {
		t.Errorf("expected '# Written', got %q", got)
	}
}

func TestOpenCode_Global_ReadSkills_Fallback(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// No canonical skills, but Claude's global skills exist
	writeTestFile(t, filepath.Join(home, ".claude", "skills", "shared-skill", "SKILL.md"), "shared skill content")

	o := &OpenCode{}
	skills, err := o.ReadSkills(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "shared-skill" {
		t.Errorf("expected fallback skills, got %v", skills)
	}
}
