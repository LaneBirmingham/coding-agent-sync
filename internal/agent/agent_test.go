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

func TestClaude_WriteSkills_InvalidName(t *testing.T) {
	root := setupTestDir(t)
	c := &Claude{}
	err := c.WriteSkills(config.Local(root), []Skill{{Name: "../escape", Content: "x"}})
	if err == nil {
		t.Fatal("expected invalid skill name error")
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

// --- Codex local tests ---

func TestCodex_ReadWriteInstructions(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, "AGENTS.md"), "# Codex instructions")

	c := &Codex{}
	inst, err := c.ReadInstructions(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if inst.Content != "# Codex instructions" {
		t.Errorf("unexpected: %q", inst.Content)
	}

	root2 := setupTestDir(t)
	err = c.WriteInstructions(config.Local(root2), inst)
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(root2, "AGENTS.md"))
	if got != "# Codex instructions" {
		t.Errorf("expected '# Codex instructions', got %q", got)
	}
}

func TestCodex_ReadInstructions_Local_FallbackPriority(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, "AGENTS.override.md"), "# Override")
	writeTestFile(t, filepath.Join(root, "AGENTS.md"), "# Agents")
	writeTestFile(t, filepath.Join(root, "TEAM_GUIDE.md"), "# Team")
	writeTestFile(t, filepath.Join(root, ".agents.md"), "# DotAgents")

	c := &Codex{}
	inst, err := c.ReadInstructions(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if inst == nil || inst.Content != "# Override" {
		t.Errorf("expected override content, got %v", inst)
	}
}

func TestCodex_ReadInstructions_Local_FallbackFiles(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, "TEAM_GUIDE.md"), "# Team")

	c := &Codex{}
	inst, err := c.ReadInstructions(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if inst == nil || inst.Content != "# Team" {
		t.Errorf("expected TEAM_GUIDE fallback content, got %v", inst)
	}
}

func TestCodex_Skills_Local_CanonicalAndFallback(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, ".codex", "skills", "legacy", "SKILL.md"), "legacy content")
	writeTestFile(t, filepath.Join(root, ".agents", "skills", "canonical", "SKILL.md"), "canonical content")

	c := &Codex{}
	skills, err := c.ReadSkills(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "canonical" {
		t.Errorf("expected canonical local skills, got %v", skills)
	}
}

func TestCodex_Skills_Local_Fallback(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, ".codex", "skills", "legacy", "SKILL.md"), "legacy content")

	c := &Codex{}
	skills, err := c.ReadSkills(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "legacy" {
		t.Errorf("expected legacy fallback skills, got %v", skills)
	}
}

func TestCodex_WriteSkills_LocalCanonical(t *testing.T) {
	root := setupTestDir(t)
	c := &Codex{}

	err := c.WriteSkills(config.Local(root), []Skill{{Name: "my-skill", Content: "skill content"}})
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(root, ".agents", "skills", "my-skill", "SKILL.md"))
	if got != "skill content" {
		t.Errorf("expected canonical local skill write, got %q", got)
	}
}

// --- Codex global tests ---

func TestCodex_Global_ReadWriteInstructions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", "")

	c := &Codex{}

	inst, err := c.ReadInstructions(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if inst != nil {
		t.Error("expected nil for missing global instructions")
	}

	err = c.WriteInstructions(config.Global(), &Instruction{Content: "# Global Codex"})
	if err != nil {
		t.Fatal(err)
	}

	inst, err = c.ReadInstructions(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if inst == nil || inst.Content != "# Global Codex" {
		t.Errorf("expected '# Global Codex', got %v", inst)
	}

	path := c.InstructionsPath(config.Global())
	if path != filepath.Join(home, ".codex", "AGENTS.md") {
		t.Errorf("unexpected global path: %q", path)
	}
}

func TestCodex_Global_Instructions_UsesCODEXHOME(t *testing.T) {
	home := t.TempDir()
	codexHome := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", codexHome)

	c := &Codex{}
	err := c.WriteInstructions(config.Global(), &Instruction{Content: "# From custom codex home"})
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(codexHome, "AGENTS.md"))
	if got != "# From custom codex home" {
		t.Errorf("expected custom CODEX_HOME write, got %q", got)
	}
}

func TestCodex_Global_Instructions_FallbackToHomeCodex(t *testing.T) {
	home := t.TempDir()
	codexHome := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", codexHome)

	writeTestFile(t, filepath.Join(home, ".codex", "AGENTS.md"), "# Home fallback")

	c := &Codex{}
	inst, err := c.ReadInstructions(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if inst == nil || inst.Content != "# Home fallback" {
		t.Errorf("expected fallback content, got %v", inst)
	}
}

func TestCodex_Global_Skills_CanonicalAndFallback(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", "")

	writeTestFile(t, filepath.Join(home, ".agents", "skills", "canonical", "SKILL.md"), "canonical content")
	writeTestFile(t, filepath.Join(home, ".codex", "skills", "legacy", "SKILL.md"), "legacy content")

	c := &Codex{}
	skills, err := c.ReadSkills(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "canonical" {
		t.Errorf("expected canonical global skills, got %v", skills)
	}
}

func TestCodex_Global_Skills_Fallback(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", "")

	writeTestFile(t, filepath.Join(home, ".codex", "skills", "legacy", "SKILL.md"), "legacy content")

	c := &Codex{}
	skills, err := c.ReadSkills(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "legacy" {
		t.Errorf("expected legacy fallback skills, got %v", skills)
	}
}

func TestCodex_WriteSkills_GlobalCanonical(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", "")

	c := &Codex{}
	err := c.WriteSkills(config.Global(), []Skill{{Name: "my-skill", Content: "global skill"}})
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(home, ".agents", "skills", "my-skill", "SKILL.md"))
	if got != "global skill" {
		t.Errorf("expected canonical global skill write, got %q", got)
	}
}

func TestCodex_Global_InstructionsPath_WindowsStyleEnv(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", `C:\Users\codex`)
	t.Setenv("CODEX_HOME", "")

	c := &Codex{}
	got := c.InstructionsPath(config.Global())
	want := filepath.Join(`C:\Users\codex`, ".codex", "AGENTS.md")
	if got != want {
		t.Errorf("expected windows-style home path %q, got %q", want, got)
	}
}

// --- Gemini local tests ---

func TestGemini_ReadWriteInstructions(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, "GEMINI.md"), "# Gemini instructions")

	g := &Gemini{}
	inst, err := g.ReadInstructions(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if inst == nil || inst.Content != "# Gemini instructions" {
		t.Errorf("expected '# Gemini instructions', got %v", inst)
	}

	root2 := setupTestDir(t)
	err = g.WriteInstructions(config.Local(root2), inst)
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(root2, "GEMINI.md"))
	if got != "# Gemini instructions" {
		t.Errorf("expected '# Gemini instructions', got %q", got)
	}
}

func TestGemini_ReadInstructions_Missing(t *testing.T) {
	root := setupTestDir(t)
	g := &Gemini{}
	inst, err := g.ReadInstructions(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if inst != nil {
		t.Error("expected nil for missing instructions")
	}
}

func TestGemini_InstructionsPath(t *testing.T) {
	g := &Gemini{}
	got := g.InstructionsPath(config.Local("/project"))
	want := filepath.Join("/project", "GEMINI.md")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestGemini_ReadWriteSkills(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, ".gemini", "skills", "my-skill", "SKILL.md"), "skill content")

	g := &Gemini{}
	skills, err := g.ReadSkills(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "my-skill" || skills[0].Content != "skill content" {
		t.Errorf("unexpected skills: %v", skills)
	}

	root2 := setupTestDir(t)
	err = g.WriteSkills(config.Local(root2), skills)
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(root2, ".gemini", "skills", "my-skill", "SKILL.md"))
	if got != "skill content" {
		t.Errorf("expected 'skill content', got %q", got)
	}
}

func TestGemini_Skills_Local_FallbackPriority(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, ".agents", "skills", "canonical", "SKILL.md"), "canonical content")
	writeTestFile(t, filepath.Join(root, ".gemini", "skills", "legacy", "SKILL.md"), "legacy content")

	g := &Gemini{}
	skills, err := g.ReadSkills(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "canonical" {
		t.Errorf("expected canonical skills, got %v", skills)
	}
}

func TestGemini_Skills_Local_Fallback(t *testing.T) {
	root := setupTestDir(t)
	writeTestFile(t, filepath.Join(root, ".gemini", "skills", "legacy", "SKILL.md"), "legacy content")

	g := &Gemini{}
	skills, err := g.ReadSkills(config.Local(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "legacy" {
		t.Errorf("expected legacy fallback skills, got %v", skills)
	}
}

// --- Gemini global tests ---

func TestGemini_Global_ReadWriteInstructions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	g := &Gemini{}

	inst, err := g.ReadInstructions(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if inst != nil {
		t.Error("expected nil for missing global instructions")
	}

	err = g.WriteInstructions(config.Global(), &Instruction{Content: "# Global Gemini"})
	if err != nil {
		t.Fatal(err)
	}

	inst, err = g.ReadInstructions(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if inst == nil || inst.Content != "# Global Gemini" {
		t.Errorf("expected '# Global Gemini', got %v", inst)
	}

	path := g.InstructionsPath(config.Global())
	if path != filepath.Join(home, ".gemini", "GEMINI.md") {
		t.Errorf("unexpected global path: %q", path)
	}
}

func TestGemini_Global_ReadWriteSkills(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	g := &Gemini{}
	writeTestFile(t, filepath.Join(home, ".gemini", "skills", "global-skill", "SKILL.md"), "global skill")

	skills, err := g.ReadSkills(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "global-skill" {
		t.Errorf("unexpected skills: %v", skills)
	}

	home2 := t.TempDir()
	t.Setenv("HOME", home2)
	err = g.WriteSkills(config.Global(), skills)
	if err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, filepath.Join(home2, ".gemini", "skills", "global-skill", "SKILL.md"))
	if got != "global skill" {
		t.Errorf("expected 'global skill', got %q", got)
	}
}

func TestGemini_Global_Skills_FallbackPriority(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeTestFile(t, filepath.Join(home, ".agents", "skills", "canonical", "SKILL.md"), "canonical content")
	writeTestFile(t, filepath.Join(home, ".gemini", "skills", "legacy", "SKILL.md"), "legacy content")

	g := &Gemini{}
	skills, err := g.ReadSkills(config.Global())
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "canonical" {
		t.Errorf("expected canonical global skills, got %v", skills)
	}
}
