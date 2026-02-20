package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
	"github.com/LaneBirmingham/coding-agent-sync/internal/sync"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	var (
		flagFrom   string
		flagTo     string
		flagDryRun bool
	)

	cmd := &cobra.Command{
		Use:   "sync [instructions|skills]",
		Short: "Sync configuration from one agent to others",
		Long:  "Sync instructions and/or skills from a source agent to one or more destination agents.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind := sync.All
			if len(args) == 1 {
				switch args[0] {
				case "instructions":
					kind = sync.Instructions
				case "skills":
					kind = sync.Skills
				default:
					return fmt.Errorf("unknown sync target %q (valid: instructions, skills)", args[0])
				}
			}
			return doSync(kind, flagFrom, flagTo, flagDryRun)
		},
	}

	cmd.Flags().StringVar(&flagFrom, "from", "", "source agent (claude, copilot, opencode)")
	cmd.Flags().StringVar(&flagTo, "to", "", "destination agent(s), comma-separated")
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "preview changes without writing")

	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")

	return cmd
}

func doSync(kind sync.ItemKind, from, to string, dryRun bool) error {
	cfg, err := buildSyncConfig(from, to, dryRun)
	if err != nil {
		return err
	}

	result, err := sync.SyncAll(cfg, kind)
	if err != nil {
		return err
	}

	for _, action := range result.Actions {
		fmt.Println(action)
	}
	return nil
}

func buildSyncConfig(fromStr, toStr string, dryRun bool) (*config.SyncConfig, error) {
	from, err := config.ParseAgent(fromStr)
	if err != nil {
		return nil, err
	}

	var targets []config.Agent
	for _, s := range strings.Split(toStr, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		a, err := config.ParseAgent(s)
		if err != nil {
			return nil, err
		}
		if a == from {
			fmt.Fprintf(os.Stderr, "warning: skipping %s (same as source)\n", a)
			continue
		}
		targets = append(targets, a)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no valid destination agents specified")
	}

	root := flagRoot
	if root == "" || root == "." {
		root, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("getting working directory: %w", err)
		}
	}

	return &config.SyncConfig{
		From:    from,
		To:      targets,
		Root:    root,
		DryRun:  dryRun,
		Verbose: flagVerbose,
	}, nil
}
