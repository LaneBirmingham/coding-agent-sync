package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// Copilot implements Agent for GitHub Copilot.
type Copilot struct{}

func (c *Copilot) Name() string { return "copilot" }

func (c *Copilot) InstructionsPath(loc config.Location) string {
	if loc.Scope == config.ScopeGlobal {
		// Copilot does not support global instructions
		return ""
	}
	return filepath.Join(loc.Root, "AGENTS.md")
}

func (c *Copilot) ReadInstructions(loc config.Location) (*Instruction, error) {
	if loc.Scope == config.ScopeGlobal {
		// Copilot does not support global instructions
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

func (c *Copilot) ReadSkills(loc config.Location) ([]Skill, error) {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining home directory: %w", err)
		}
		return readSkillsFromDir(filepath.Join(home, ".copilot", "skills"))
	}
	return readSkillsFromDir(filepath.Join(loc.Root, ".github", "skills"))
}

func (c *Copilot) WriteInstructions(loc config.Location, inst *Instruction) error {
	if loc.Scope == config.ScopeGlobal {
		return fmt.Errorf("copilot does not support global instructions")
	}
	return writeFile(filepath.Join(loc.Root, "AGENTS.md"), inst.Content)
}

func (c *Copilot) WriteSkills(loc config.Location, skills []Skill) error {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("determining home directory: %w", err)
		}
		return writeSkillsToDir(filepath.Join(home, ".copilot", "skills"), skills)
	}
	return writeSkillsToDir(filepath.Join(loc.Root, ".github", "skills"), skills)
}
