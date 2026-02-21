package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// OpenCode implements Agent for OpenCode.
type OpenCode struct{}

func (o *OpenCode) Name() string { return "opencode" }

func (o *OpenCode) InstructionsPath(loc config.Location) string {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, ".config", "opencode", "AGENTS.md")
	}
	return filepath.Join(loc.Root, "AGENTS.md")
}

func (o *OpenCode) ReadInstructions(loc config.Location) (*Instruction, error) {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining home directory: %w", err)
		}
		// Try canonical path first
		content, err := readFile(filepath.Join(home, ".config", "opencode", "AGENTS.md"))
		if err != nil {
			return nil, err
		}
		if content != "" {
			return &Instruction{Content: content}, nil
		}
		// Fallback: read from Claude's global config
		content, err = readFile(filepath.Join(home, ".claude", "CLAUDE.md"))
		if err != nil {
			return nil, err
		}
		if content != "" {
			return &Instruction{Content: content}, nil
		}
		return nil, nil
	}
	content, err := readFile(filepath.Join(loc.Root, "AGENTS.md"))
	if err != nil {
		return nil, err
	}
	if content == "" {
		return nil, nil
	}
	return &Instruction{Content: content}, nil
}

func (o *OpenCode) ReadSkills(loc config.Location) ([]Skill, error) {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining home directory: %w", err)
		}
		// Try canonical path first
		skills, err := readSkillsFromDir(filepath.Join(home, ".config", "opencode", "skills"))
		if err != nil {
			return nil, err
		}
		if len(skills) > 0 {
			return skills, nil
		}
		// Fallback: read from Claude's global skills
		return readSkillsFromDir(filepath.Join(home, ".claude", "skills"))
	}
	return readSkillsFromDir(filepath.Join(loc.Root, ".opencode", "skills"))
}

func (o *OpenCode) WriteInstructions(loc config.Location, inst *Instruction) error {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("determining home directory: %w", err)
		}
		return writeFile(filepath.Join(home, ".config", "opencode", "AGENTS.md"), inst.Content)
	}
	return writeFile(filepath.Join(loc.Root, "AGENTS.md"), inst.Content)
}

func (o *OpenCode) WriteSkills(loc config.Location, skills []Skill) error {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("determining home directory: %w", err)
		}
		return writeSkillsToDir(filepath.Join(home, ".config", "opencode", "skills"), skills)
	}
	return writeSkillsToDir(filepath.Join(loc.Root, ".opencode", "skills"), skills)
}
