package agent

import (
	"os"
	"path/filepath"
	"testing"
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

// --- Claude tests ---

func TestClaude_ReadInstructions_RootFile(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, "CLAUDE.md"), "# Instructions")

	c := &Claude{}
	inst, err := c.ReadInstructions(root)
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
	inst, err := c.ReadInstructions(root)
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
	inst, err := c.ReadInstructions(root)
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
	err := c.WriteInstructions(root, &Instruction{Content: "# Test"})
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
	got := c.InstructionsPath("/project")
	want := filepath.Join("/project", "CLAUDE.md")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestClaude_ReadWriteSkills(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, ".claude", "skills", "my-skill", "SKILL.md"), "skill content")

	c := &Claude{}
	skills, err := c.ReadSkills(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "my-skill" || skills[0].Content != "skill content" {
		t.Errorf("unexpected skills: %v", skills)
	}

	root2 := setupTestDir(t)
	err = c.WriteSkills(root2, skills)
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(root2, ".claude", "skills", "my-skill", "SKILL.md"))
	if got != "skill content" {
		t.Errorf("expected 'skill content', got %q", got)
	}
}

// --- Copilot tests ---

func TestCopilot_ReadWriteInstructions(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, "AGENTS.md"), "# Copilot instructions")

	c := &Copilot{}
	inst, err := c.ReadInstructions(root)
	if err != nil {
		t.Fatal(err)
	}
	if inst.Content != "# Copilot instructions" {
		t.Errorf("unexpected: %q", inst.Content)
	}

	root2 := setupTestDir(t)
	err = c.WriteInstructions(root2, inst)
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
	skills, err := c.ReadSkills(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "s1" {
		t.Errorf("unexpected skills: %v", skills)
	}
}

// --- OpenCode tests ---

func TestOpenCode_ReadWriteInstructions(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, "AGENTS.md"), "# OC instructions")

	o := &OpenCode{}
	inst, err := o.ReadInstructions(root)
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
	skills, err := o.ReadSkills(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "s1" {
		t.Errorf("unexpected skills: %v", skills)
	}
}
