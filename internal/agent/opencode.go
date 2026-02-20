package agent

import (
	"path/filepath"
)

// OpenCode implements Agent for OpenCode.
type OpenCode struct{}

func (o *OpenCode) Name() string { return "opencode" }

func (o *OpenCode) InstructionsPath(root string) string {
	return filepath.Join(root, "AGENTS.md")
}

func (o *OpenCode) ReadInstructions(root string) (*Instruction, error) {
	content, err := readFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		return nil, err
	}
	if content == "" {
		return nil, nil
	}
	return &Instruction{Content: content}, nil
}

func (o *OpenCode) ReadSkills(root string) ([]Skill, error) {
	return readSkillsFromDir(filepath.Join(root, ".opencode", "skills"))
}

func (o *OpenCode) WriteInstructions(root string, inst *Instruction) error {
	return writeFile(filepath.Join(root, "AGENTS.md"), inst.Content)
}

func (o *OpenCode) WriteSkills(root string, skills []Skill) error {
	return writeSkillsToDir(filepath.Join(root, ".opencode", "skills"), skills)
}
