package sync

import (
	"fmt"

	"github.com/LaneBirmingham/coding-agent-sync/internal/agent"
	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
)

// SyncInstructions syncs instructions from source to destination.
func SyncInstructions(root string, from, to config.Agent, dryRun bool) (SyncAction, error) {
	src, err := agent.Get(from)
	if err != nil {
		return SyncAction{}, err
	}
	dst, err := agent.Get(to)
	if err != nil {
		return SyncAction{}, err
	}

	action := SyncAction{Kind: Instructions, From: from, To: to}

	// Detect shared file path (e.g., both use AGENTS.md)
	if src.InstructionsPath(root) == dst.InstructionsPath(root) {
		action.Status = "noop"
		action.Detail = fmt.Sprintf("already in sync (both use %s)", src.InstructionsPath(root))
		return action, nil
	}

	inst, err := src.ReadInstructions(root)
	if err != nil {
		return SyncAction{}, fmt.Errorf("reading instructions from %s: %w", from, err)
	}
	if inst == nil {
		action.Status = "skipped"
		action.Detail = "skipped (no source file found)"
		return action, nil
	}

	if dryRun {
		action.Status = "dry-run"
		action.Detail = fmt.Sprintf("would write (%d bytes)", len(inst.Content))
		return action, nil
	}

	if err := dst.WriteInstructions(root, inst); err != nil {
		return SyncAction{}, fmt.Errorf("writing instructions to %s: %w", to, err)
	}

	action.Status = "synced"
	action.Detail = fmt.Sprintf("synced (%d bytes)", len(inst.Content))
	return action, nil
}
