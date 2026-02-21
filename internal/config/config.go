package config

import (
	"fmt"
	"strings"
)

// Agent represents a supported coding agent.
type Agent string

const (
	Claude   Agent = "claude"
	Copilot  Agent = "copilot"
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

// Scope represents whether config is project-level or user-level.
type Scope string

const (
	ScopeLocal  Scope = "local"
	ScopeGlobal Scope = "global"
)

// ParseScope converts a string to a Scope, returning an error if invalid.
func ParseScope(s string) (Scope, error) {
	switch Scope(strings.ToLower(s)) {
	case ScopeLocal, ScopeGlobal:
		return Scope(strings.ToLower(s)), nil
	default:
		return "", fmt.Errorf("unknown scope %q (valid: local, global)", s)
	}
}

// Location specifies where to read/write agent config.
type Location struct {
	Root  string // Project root directory (used for ScopeLocal)
	Scope Scope
}

// Local returns a Location for project-level config at the given root.
func Local(root string) Location { return Location{Root: root, Scope: ScopeLocal} }

// Global returns a Location for user-level config.
func Global() Location { return Location{Scope: ScopeGlobal} }

// SyncConfig holds the configuration for a sync operation.
type SyncConfig struct {
	From      Agent
	To        []Agent
	Root      string
	FromScope Scope
	ToScope   Scope
	DryRun    bool
	Verbose   bool
}
