package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// Gemini implements Agent for Gemini CLI.
type Gemini struct{}

func (g *Gemini) Name() string { return "gemini" }

func (g *Gemini) InstructionsPath(loc config.Location) string {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, ".gemini", "GEMINI.md")
	}
	return filepath.Join(loc.Root, "GEMINI.md")
}

func (g *Gemini) ReadInstructions(loc config.Location) (*Instruction, error) {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining home directory: %w", err)
		}
		content, err := readFile(filepath.Join(home, ".gemini", "GEMINI.md"))
		if err != nil {
			return nil, err
		}
		if content == "" {
			return nil, nil
		}
		return &Instruction{Content: content}, nil
	}

	content, err := readFile(filepath.Join(loc.Root, "GEMINI.md"))
	if err != nil {
		return nil, err
	}
	if content == "" {
		return nil, nil
	}
	return &Instruction{Content: content}, nil
}

func (g *Gemini) ReadSkills(loc config.Location) ([]Skill, error) {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining home directory: %w", err)
		}
		// .agents/skills/ takes precedence over .gemini/skills/
		skills, err := readSkillsFromDir(filepath.Join(home, ".agents", "skills"))
		if err != nil {
			return nil, err
		}
		if len(skills) > 0 {
			return skills, nil
		}
		return readSkillsFromDir(filepath.Join(home, ".gemini", "skills"))
	}

	// .agents/skills/ takes precedence over .gemini/skills/
	skills, err := readSkillsFromDir(filepath.Join(loc.Root, ".agents", "skills"))
	if err != nil {
		return nil, err
	}
	if len(skills) > 0 {
		return skills, nil
	}
	return readSkillsFromDir(filepath.Join(loc.Root, ".gemini", "skills"))
}

func (g *Gemini) WriteInstructions(loc config.Location, inst *Instruction) error {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("determining home directory: %w", err)
		}
		return writeFile(filepath.Join(home, ".gemini", "GEMINI.md"), inst.Content)
	}
	return writeFile(filepath.Join(loc.Root, "GEMINI.md"), inst.Content)
}

func (g *Gemini) WriteSkills(loc config.Location, skills []Skill) error {
	if loc.Scope == config.ScopeGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("determining home directory: %w", err)
		}
		return writeSkillsToDir(filepath.Join(home, ".gemini", "skills"), skills)
	}
	return writeSkillsToDir(filepath.Join(loc.Root, ".gemini", "skills"), skills)
}
