package sync

import (
	"fmt"

	"github.com/LaneBirmingham/coding-agent-sync/internal/agent"
	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// SyncInstructions syncs instructions from source to destination.
func SyncInstructions(cfg *config.SyncConfig, to config.Agent) (SyncAction, error) {
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
		Kind:      Instructions,
		From:      cfg.From,
		To:        to,
		FromScope: cfg.FromScope,
		ToScope:   cfg.ToScope,
	}

	// Check if destination supports instructions at this scope
	dstPath := dst.InstructionsPath(dstLoc)
	if dstPath == "" {
		action.Status = "skipped"
		action.Detail = fmt.Sprintf("skipped (%s does not support %s instructions)", to, dstLoc.Scope)
		return action, nil
	}

	// Check if source supports instructions at this scope
	srcPath := src.InstructionsPath(srcLoc)
	if srcPath == "" {
		action.Status = "skipped"
		action.Detail = fmt.Sprintf("skipped (%s does not support %s instructions)", cfg.From, srcLoc.Scope)
		return action, nil
	}

	// Detect shared file path (e.g., both use AGENTS.md at same scope)
	if srcPath == dstPath {
		action.Status = "noop"
		action.Detail = fmt.Sprintf("already in sync (both use %s)", srcPath)
		return action, nil
	}

	inst, err := src.ReadInstructions(srcLoc)
	if err != nil {
		return SyncAction{}, fmt.Errorf("reading instructions from %s: %w", cfg.From, err)
	}
	if inst == nil {
		action.Status = "skipped"
		action.Detail = "skipped (no source file found)"
		return action, nil
	}

	if cfg.DryRun {
		action.Status = "dry-run"
		action.Detail = fmt.Sprintf("would write (%d bytes)", len(inst.Content))
		return action, nil
	}

	if err := dst.WriteInstructions(dstLoc, inst); err != nil {
		return SyncAction{}, fmt.Errorf("writing instructions to %s: %w", to, err)
	}

	action.Status = "synced"
	action.Detail = fmt.Sprintf("synced (%d bytes)", len(inst.Content))
	return action, nil
}
