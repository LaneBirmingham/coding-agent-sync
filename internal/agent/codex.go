package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// Codex implements Agent for OpenAI Codex.
type Codex struct{}

func (c *Codex) Name() string { return "codex" }

func (c *Codex) InstructionsPath(loc config.Location) string {
	if loc.Scope == config.ScopeGlobal {
		codexHome, err := resolveCodexHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(codexHome, "AGENTS.md")
	}
	return filepath.Join(loc.Root, "AGENTS.md")
}

func (c *Codex) ReadInstructions(loc config.Location) (*Instruction, error) {
	if loc.Scope == config.ScopeGlobal {
		codexHome, home, err := resolveCodexAndHomeDirs()
		if err != nil {
			return nil, fmt.Errorf("determining codex home directory: %w", err)
		}
		paths := uniquePaths(
			filepath.Join(codexHome, "AGENTS.override.md"),
			filepath.Join(codexHome, "AGENTS.md"),
			filepath.Join(home, ".codex", "AGENTS.override.md"),
			filepath.Join(home, ".codex", "AGENTS.md"),
		)
		return readFirstInstruction(paths)
	}

	paths := []string{
		filepath.Join(loc.Root, "AGENTS.override.md"),
		filepath.Join(loc.Root, "AGENTS.md"),
		filepath.Join(loc.Root, "TEAM_GUIDE.md"),
		filepath.Join(loc.Root, ".agents.md"),
	}
	return readFirstInstruction(paths)
}

func (c *Codex) ReadSkills(loc config.Location) ([]Skill, error) {
	if loc.Scope == config.ScopeGlobal {
		home, err := resolveHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining home directory: %w", err)
		}
		skills, err := readSkillsFromDir(filepath.Join(home, ".agents", "skills"))
		if err != nil {
			return nil, err
		}
		if len(skills) > 0 {
			return skills, nil
		}

		codexHome, err := resolveCodexHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining codex home directory: %w", err)
		}
		return readSkillsFromDir(filepath.Join(codexHome, "skills"))
	}

	skills, err := readSkillsFromDir(filepath.Join(loc.Root, ".agents", "skills"))
	if err != nil {
		return nil, err
	}
	if len(skills) > 0 {
		return skills, nil
	}
	return readSkillsFromDir(filepath.Join(loc.Root, ".codex", "skills"))
}

func (c *Codex) WriteInstructions(loc config.Location, inst *Instruction) error {
	if loc.Scope == config.ScopeGlobal {
		codexHome, err := resolveCodexHomeDir()
		if err != nil {
			return fmt.Errorf("determining codex home directory: %w", err)
		}
		return writeFile(filepath.Join(codexHome, "AGENTS.md"), inst.Content)
	}
	return writeFile(filepath.Join(loc.Root, "AGENTS.md"), inst.Content)
}

func (c *Codex) WriteSkills(loc config.Location, skills []Skill) error {
	if loc.Scope == config.ScopeGlobal {
		home, err := resolveHomeDir()
		if err != nil {
			return fmt.Errorf("determining home directory: %w", err)
		}
		return writeSkillsToDir(filepath.Join(home, ".agents", "skills"), skills)
	}
	return writeSkillsToDir(filepath.Join(loc.Root, ".agents", "skills"), skills)
}

func readFirstInstruction(paths []string) (*Instruction, error) {
	for _, path := range paths {
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

func uniquePaths(paths ...string) []string {
	seen := make(map[string]struct{}, len(paths))
	var out []string
	for _, p := range paths {
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

func resolveCodexAndHomeDirs() (string, string, error) {
	home, err := resolveHomeDir()
	if err != nil {
		return "", "", err
	}
	codexHome := os.Getenv("CODEX_HOME")
	if codexHome == "" {
		codexHome = filepath.Join(home, ".codex")
	}
	return codexHome, home, nil
}

func resolveCodexHomeDir() (string, error) {
	codexHome, _, err := resolveCodexAndHomeDirs()
	return codexHome, err
}

func resolveHomeDir() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		return userProfile, nil
	}
	homeDrive := os.Getenv("HOMEDRIVE")
	homePath := os.Getenv("HOMEPATH")
	if homeDrive != "" && homePath != "" {
		return homeDrive + homePath, nil
	}
	return os.UserHomeDir()
}
