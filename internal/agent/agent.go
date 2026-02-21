package agent

import (
	"fmt"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// Instruction represents the main agent instruction file content.
type Instruction struct {
	Content string
}

// Skill represents a single skill definition with optional frontmatter.
type Skill struct {
	Name    string // Directory name (e.g., "my-skill")
	Content string // Full SKILL.md content including frontmatter
}

// Agent defines the interface for reading and writing agent configuration.
type Agent interface {
	Name() string
	InstructionsPath(loc config.Location) string
	ReadInstructions(loc config.Location) (*Instruction, error)
	ReadSkills(loc config.Location) ([]Skill, error)
	WriteInstructions(loc config.Location, inst *Instruction) error
	WriteSkills(loc config.Location, skills []Skill) error
}

// registry maps agent types to their implementations.
var registry = map[config.Agent]Agent{
	config.Claude:   &Claude{},
	config.Copilot:  &Copilot{},
	config.OpenCode: &OpenCode{},
}

// Get returns the Agent implementation for the given agent type.
func Get(a config.Agent) (Agent, error) {
	impl, ok := registry[a]
	if !ok {
		return nil, fmt.Errorf("no implementation for agent %q", a)
	}
	return impl, nil
}
