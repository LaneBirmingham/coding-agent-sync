package sync

import (
	"fmt"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// ItemKind controls which items to sync.
type ItemKind int

const (
	All          ItemKind = iota
	Instructions
	Skills
)

// SyncAction represents the outcome of a single sync operation.
type SyncAction struct {
	Kind      ItemKind
	From      config.Agent
	To        config.Agent
	FromScope config.Scope
	ToScope   config.Scope
	Status    string // "synced", "skipped", "dry-run", "noop"
	Detail    string
}

func (a SyncAction) String() string {
	label := "instructions"
	if a.Kind == Skills {
		label = "skills"
	}
	scope := ""
	if a.FromScope == config.ScopeGlobal || a.ToScope == config.ScopeGlobal {
		scope = fmt.Sprintf(" [%s→%s]", a.FromScope, a.ToScope)
	}
	return fmt.Sprintf("%s: %s → %s%s: %s", label, a.From, a.To, scope, a.Detail)
}

// Result holds the output from a sync operation.
type Result struct {
	Actions []SyncAction
}

// ArchiveAction represents the outcome of a single export/import operation.
type ArchiveAction struct {
	Kind   ItemKind
	Agent  config.Agent
	Scope  config.Scope
	Status string // "exported", "imported", "skipped", "dry-run"
	Detail string
}

func (a ArchiveAction) String() string {
	label := "instructions"
	if a.Kind == Skills {
		label = "skills"
	}
	scope := ""
	if a.Scope == config.ScopeGlobal {
		scope = fmt.Sprintf(" [%s]", a.Scope)
	}
	return fmt.Sprintf("%s: %s%s: %s", label, a.Agent, scope, a.Detail)
}

// ArchiveResult holds the output from an export or import operation.
type ArchiveResult struct {
	Actions  []ArchiveAction
	Warnings []string
}

// SyncAll runs the sync operation for each target agent.
func SyncAll(cfg *config.SyncConfig, kind ItemKind) (*Result, error) {
	var result Result

	for _, to := range cfg.To {
		if kind == All || kind == Instructions {
			action, err := SyncInstructions(cfg, to)
			if err != nil {
				return nil, err
			}
			result.Actions = append(result.Actions, action)
		}

		if kind == All || kind == Skills {
			action, err := SyncSkills(cfg, to)
			if err != nil {
				return nil, err
			}
			result.Actions = append(result.Actions, action)
		}
	}

	return &result, nil
}
