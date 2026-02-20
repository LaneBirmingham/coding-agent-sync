package agent

import (
	"path/filepath"
)

// Claude implements Agent for Claude Code.
type Claude struct{}

func (c *Claude) Name() string { return "claude" }

func (c *Claude) InstructionsPath(root string) string {
	return filepath.Join(root, "CLAUDE.md")
}

func (c *Claude) ReadInstructions(root string) (*Instruction, error) {
	// Try .claude/CLAUDE.md first, then root CLAUDE.md
	for _, path := range []string{
		filepath.Join(root, ".claude", "CLAUDE.md"),
		filepath.Join(root, "CLAUDE.md"),
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

func (c *Claude) ReadSkills(root string) ([]Skill, error) {
	return readSkillsFromDir(filepath.Join(root, ".claude", "skills"))
}

func (c *Claude) WriteInstructions(root string, inst *Instruction) error {
	return writeFile(filepath.Join(root, "CLAUDE.md"), inst.Content)
}

func (c *Claude) WriteSkills(root string, skills []Skill) error {
	return writeSkillsToDir(filepath.Join(root, ".claude", "skills"), skills)
}
