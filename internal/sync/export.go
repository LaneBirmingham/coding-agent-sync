package sync

import (
	"fmt"
	"time"

	"github.com/LaneBirmingham/coding-agent-sync/internal/agent"
	"github.com/LaneBirmingham/coding-agent-sync/internal/archive"
	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// Export reads config from an agent and writes it to a ZIP archive.
func Export(cfg *config.ExportConfig) (*ArchiveResult, error) {
	src, err := agent.Get(cfg.From)
	if err != nil {
		return nil, err
	}

	loc := config.Location{Root: cfg.Root, Scope: cfg.Scope}
	result := &ArchiveResult{}

	// Read instructions
	inst, err := src.ReadInstructions(loc)
	if err != nil {
		return nil, fmt.Errorf("reading instructions from %s: %w", cfg.From, err)
	}

	instAction := ArchiveAction{
		Kind:  Instructions,
		Agent: cfg.From,
		Scope: cfg.Scope,
	}
	if inst == nil || inst.Content == "" {
		instAction.Status = "skipped"
		instAction.Detail = "skipped (no instructions found)"
	} else if cfg.DryRun {
		instAction.Status = "dry-run"
		instAction.Detail = fmt.Sprintf("would export (%d bytes)", len(inst.Content))
	} else {
		instAction.Status = "exported"
		instAction.Detail = fmt.Sprintf("exported (%d bytes)", len(inst.Content))
	}
	result.Actions = append(result.Actions, instAction)

	// Read skills
	skills, err := src.ReadSkills(loc)
	if err != nil {
		return nil, fmt.Errorf("reading skills from %s: %w", cfg.From, err)
	}

	skillAction := ArchiveAction{
		Kind:  Skills,
		Agent: cfg.From,
		Scope: cfg.Scope,
	}
	if len(skills) == 0 {
		skillAction.Status = "skipped"
		skillAction.Detail = "skipped (no skills found)"
	} else if cfg.DryRun {
		skillAction.Status = "dry-run"
		skillAction.Detail = fmt.Sprintf("would export %d skill(s): %s", len(skills), skillNames(skills))
	} else {
		skillAction.Status = "exported"
		skillAction.Detail = fmt.Sprintf("exported %d skill(s): %s", len(skills), skillNames(skills))
	}
	result.Actions = append(result.Actions, skillAction)

	if cfg.DryRun {
		return result, nil
	}

	// Build and write archive
	a := &archive.Archive{
		Manifest: &archive.Manifest{
			Version:    archive.FormatVersion,
			Agent:      string(cfg.From),
			Scope:      string(cfg.Scope),
			ExportedAt: time.Now().UTC(),
			CASVersion: cfg.CASVersion,
		},
	}
	if inst != nil && inst.Content != "" {
		a.Instructions = inst
	}
	a.Skills = skills

	if err := archive.Write(cfg.Output, a); err != nil {
		return nil, fmt.Errorf("writing archive: %w", err)
	}

	return result, nil
}
