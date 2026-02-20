package config

import (
	"fmt"
	"strings"
)

// Agent represents a supported coding agent.
type Agent string

const (
	Claude  Agent = "claude"
	Copilot Agent = "copilot"
	OpenCode Agent = "opencode"
)

// ValidAgents lists all supported agents.
var ValidAgents = []Agent{Claude, Copilot, OpenCode}

// ParseAgent converts a string to an Agent, returning an error if invalid.
func ParseAgent(s string) (Agent, error) {
	switch Agent(strings.ToLower(s)) {
	case Claude, Copilot, OpenCode:
		return Agent(strings.ToLower(s)), nil
	default:
		return "", fmt.Errorf("unknown agent %q (valid: claude, copilot, opencode)", s)
	}
}

// SyncConfig holds the configuration for a sync operation.
type SyncConfig struct {
	From    Agent
	To      []Agent
	Root    string
	DryRun  bool
	Verbose bool
}
