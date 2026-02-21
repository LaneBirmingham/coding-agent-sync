package sync

import (
	"fmt"
	"strings"

	"github.com/LaneBirmingham/coding-agent-sync/internal/agent"
	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// SyncSkills syncs skills from source to destination.
func SyncSkills(cfg *config.SyncConfig, to config.Agent) (SyncAction, error) {
	src, err := agent.Get(cfg.From)
	if err != nil {
		return SyncAction{}, err
	}
	dst, err := agent.Get(to)
	if err != nil {
		return SyncAction{}, err
	}

	srcLoc := config.Location{Root: cfg.Root, Scope: cfg.FromScope}
	dstLoc := config.Location{Root: cfg.Root, Scope: cfg.ToScope}

	action := SyncAction{
		Kind:      Skills,
		From:      cfg.From,
		To:        to,
		FromScope: cfg.FromScope,
		ToScope:   cfg.ToScope,
	}

	skills, err := src.ReadSkills(srcLoc)
	if err != nil {
		return SyncAction{}, fmt.Errorf("reading skills from %s: %w", cfg.From, err)
	}
	if len(skills) == 0 {
		action.Status = "skipped"
		action.Detail = "skipped (no skills found)"
		return action, nil
	}

	if cfg.DryRun {
		action.Status = "dry-run"
		action.Detail = fmt.Sprintf("would write %d skill(s): %s", len(skills), skillNames(skills))
		return action, nil
	}

	if err := dst.WriteSkills(dstLoc, skills); err != nil {
		return SyncAction{}, fmt.Errorf("writing skills to %s: %w", to, err)
	}

	action.Status = "synced"
	action.Detail = fmt.Sprintf("synced %d skill(s): %s", len(skills), skillNames(skills))
	return action, nil
}

func skillNames(skills []agent.Skill) string {
	names := make([]string, len(skills))
	for i, s := range skills {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}
