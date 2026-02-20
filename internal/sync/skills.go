package sync

import (
	"fmt"

	"github.com/LaneBirmingham/coding-agent-sync/internal/agent"
	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// SyncSkills syncs skills from source to destination.
func SyncSkills(root string, from, to config.Agent, dryRun bool) (SyncAction, error) {
	src, err := agent.Get(from)
	if err != nil {
		return SyncAction{}, err
	}
	dst, err := agent.Get(to)
	if err != nil {
		return SyncAction{}, err
	}

	action := SyncAction{Kind: Skills, From: from, To: to}

	skills, err := src.ReadSkills(root)
	if err != nil {
		return SyncAction{}, fmt.Errorf("reading skills from %s: %w", from, err)
	}
	if len(skills) == 0 {
		action.Status = "skipped"
		action.Detail = "skipped (no skills found)"
		return action, nil
	}

	if dryRun {
		action.Status = "dry-run"
		action.Detail = fmt.Sprintf("would write %d skill(s)", len(skills))
		return action, nil
	}

	if err := dst.WriteSkills(root, skills); err != nil {
		return SyncAction{}, fmt.Errorf("writing skills to %s: %w", to, err)
	}

	action.Status = "synced"
	action.Detail = fmt.Sprintf("synced %d skill(s)", len(skills))
	return action, nil
}
