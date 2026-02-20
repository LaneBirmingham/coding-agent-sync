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
	Kind   ItemKind
	From   config.Agent
	To     config.Agent
	Status string // "synced", "skipped", "dry-run", "noop"
	Detail string
}

func (a SyncAction) String() string {
	label := "instructions"
	if a.Kind == Skills {
		label = "skills"
	}
	return fmt.Sprintf("%s: %s → %s: %s", label, a.From, a.To, a.Detail)
}

// Result holds the output from a sync operation.
type Result struct {
	Actions []SyncAction
}

// SyncAll runs the sync operation for each target agent.
func SyncAll(cfg *config.SyncConfig, kind ItemKind) (*Result, error) {
	var result Result

	for _, to := range cfg.To {
		if kind == All || kind == Instructions {
			action, err := SyncInstructions(cfg.Root, cfg.From, to, cfg.DryRun)
			if err != nil {
				return nil, err
			}
			result.Actions = append(result.Actions, action)
		}

		if kind == All || kind == Skills {
			action, err := SyncSkills(cfg.Root, cfg.From, to, cfg.DryRun)
			if err != nil {
				return nil, err
			}
			result.Actions = append(result.Actions, action)
		}
	}

	return &result, nil
}
