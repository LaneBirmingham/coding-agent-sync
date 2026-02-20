package agent

import (
	"path/filepath"
)

// Copilot implements Agent for GitHub Copilot.
type Copilot struct{}

func (c *Copilot) Name() string { return "copilot" }

func (c *Copilot) InstructionsPath(root string) string {
	return filepath.Join(root, "AGENTS.md")
}

func (c *Copilot) ReadInstructions(root string) (*Instruction, error) {
	content, err := readFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		return nil, err
	}
	if content == "" {
		return nil, nil
	}
	return &Instruction{Content: content}, nil
}

func (c *Copilot) ReadSkills(root string) ([]Skill, error) {
	return readSkillsFromDir(filepath.Join(root, ".github", "skills"))
}

func (c *Copilot) WriteInstructions(root string, inst *Instruction) error {
	return writeFile(filepath.Join(root, "AGENTS.md"), inst.Content)
}

func (c *Copilot) WriteSkills(root string, skills []Skill) error {
	return writeSkillsToDir(filepath.Join(root, ".github", "skills"), skills)
}
