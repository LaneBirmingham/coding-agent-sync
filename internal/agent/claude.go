package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// Claude implements Agent for Claude Code.
type Claude struct{}

func (c *Claude) Name() string { return "claude" }

func (c *Claude) InstructionsPath(loc config.Location) string {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, ".claude", "CLAUDE.md")
	}
	return filepath.Join(loc.Root, "CLAUDE.md")
}

func (c *Claude) ReadInstructions(loc config.Location) (*Instruction, error) {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining home directory: %w", err)
		}
		content, err := readFile(filepath.Join(home, ".claude", "CLAUDE.md"))
		if err != nil {
			return nil, err
		}
		if content == "" {
			return nil, nil
		}
		return &Instruction{Content: content}, nil
	}

	// Try .claude/CLAUDE.md first, then root CLAUDE.md
	for _, path := range []string{
		filepath.Join(loc.Root, ".claude", "CLAUDE.md"),
		filepath.Join(loc.Root, "CLAUDE.md"),
	} {
		content, err := readFile(path)
		if err != nil {
			return nil, err
		}
		if content != "" {
			return &Instruction{Content: content}, nil
		}
	}
	return nil, nil
}

func (c *Claude) ReadSkills(loc config.Location) ([]Skill, error) {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining home directory: %w", err)
		}
		return readSkillsFromDir(filepath.Join(home, ".claude", "skills"))
	}
	return readSkillsFromDir(filepath.Join(loc.Root, ".claude", "skills"))
}

func (c *Claude) WriteInstructions(loc config.Location, inst *Instruction) error {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("determining home directory: %w", err)
		}
		return writeFile(filepath.Join(home, ".claude", "CLAUDE.md"), inst.Content)
	}
	return writeFile(filepath.Join(loc.Root, "CLAUDE.md"), inst.Content)
}

func (c *Claude) WriteSkills(loc config.Location, skills []Skill) error {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("determining home directory: %w", err)
		}
		return writeSkillsToDir(filepath.Join(home, ".claude", "skills"), skills)
	}
	return writeSkillsToDir(filepath.Join(loc.Root, ".claude", "skills"), skills)
}
