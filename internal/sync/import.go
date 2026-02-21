package sync

import (
	"fmt"

	"github.com/LaneBirmingham/coding-agent-sync/internal/agent"
	"github.com/LaneBirmingham/coding-agent-sync/internal/archive"
	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// Import reads a ZIP archive and writes its contents to one or more agents.
func Import(cfg *config.ImportConfig) (*ArchiveResult, error) {
	a, err := archive.Read(cfg.Input)
	if err != nil {
		return nil, fmt.Errorf("reading archive: %w", err)
	}

	result := &ArchiveResult{}

	for _, to := range cfg.To {
		// Warn on agent mismatch
		if config.Agent(a.Manifest.Agent) != to {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("archive was exported from %s, importing to %s", a.Manifest.Agent, to))
		}
		// Warn on scope mismatch
		if config.Scope(a.Manifest.Scope) != cfg.Scope {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("archive scope is %s, importing to %s scope", a.Manifest.Scope, cfg.Scope))
		}

		dst, err := agent.Get(to)
		if err != nil {
			return nil, err
		}

		loc := config.Location{Root: cfg.Root, Scope: cfg.Scope}

		// Import instructions
		instAction := ArchiveAction{
			Kind:  Instructions,
			Agent: to,
			Scope: cfg.Scope,
		}
		dstPath := dst.InstructionsPath(loc)
		if dstPath == "" {
			instAction.Status = "skipped"
			instAction.Detail = fmt.Sprintf("skipped (%s does not support %s instructions)", to, cfg.Scope)
		} else if a.Instructions == nil || a.Instructions.Content == "" {
			instAction.Status = "skipped"
			instAction.Detail = "skipped (no instructions in archive)"
		} else if cfg.DryRun {
			instAction.Status = "dry-run"
			instAction.Detail = fmt.Sprintf("would import (%d bytes)", len(a.Instructions.Content))
		} else {
			if err := dst.WriteInstructions(loc, a.Instructions); err != nil {
				return nil, fmt.Errorf("writing instructions to %s: %w", to, err)
			}
			instAction.Status = "imported"
			instAction.Detail = fmt.Sprintf("imported (%d bytes)", len(a.Instructions.Content))
		}
		result.Actions = append(result.Actions, instAction)

		// Import skills
		skillAction := ArchiveAction{
			Kind:  Skills,
			Agent: to,
			Scope: cfg.Scope,
		}
		if len(a.Skills) == 0 {
			skillAction.Status = "skipped"
			skillAction.Detail = "skipped (no skills in archive)"
		} else if cfg.DryRun {
			skillAction.Status = "dry-run"
			skillAction.Detail = fmt.Sprintf("would import %d skill(s): %s", len(a.Skills), skillNames(a.Skills))
		} else {
			if err := dst.WriteSkills(loc, a.Skills); err != nil {
				return nil, fmt.Errorf("writing skills to %s: %w", to, err)
			}
			skillAction.Status = "imported"
			skillAction.Detail = fmt.Sprintf("imported %d skill(s): %s", len(a.Skills), skillNames(a.Skills))
		}
		result.Actions = append(result.Actions, skillAction)
	}

	return result, nil
}
