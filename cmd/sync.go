package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
	"github.com/LaneBirmingham/coding-agent-sync/internal/sync"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	var (
		flagFrom      string
		flagTo        string
		flagDryRun    bool
		flagScope     string
		flagFromScope string
		flagToScope   string
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
			return doSync(kind, flagFrom, flagTo, flagDryRun, flagScope, flagFromScope, flagToScope)
		},
	}

	cmd.Flags().StringVar(&flagFrom, "from", "", "source agent (claude, copilot, opencode)")
	cmd.Flags().StringVar(&flagTo, "to", "", "destination agent(s), comma-separated")
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "preview changes without writing")
	cmd.Flags().StringVar(&flagScope, "scope", "", "set both from and to scope (local, global)")
	cmd.Flags().StringVar(&flagFromScope, "from-scope", "", "source scope (overrides --scope)")
	cmd.Flags().StringVar(&flagToScope, "to-scope", "", "destination scope (overrides --scope)")

	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")

	return cmd
}

func doSync(kind sync.ItemKind, from, to string, dryRun bool, scope, fromScope, toScope string) error {
	cfg, err := buildSyncConfig(from, to, dryRun, scope, fromScope, toScope)
	if err != nil {
		return err
	}

	if cfg.Verbose {
		targets := make([]string, 0, len(cfg.To))
		for _, t := range cfg.To {
			targets = append(targets, string(t))
		}
		slices.Sort(targets)
		fmt.Fprintf(os.Stderr, "verbose: sync kind=%s from=%s(%s) to=%s(%s) root=%s dry-run=%t\n",
			itemKindLabel(kind), cfg.From, cfg.FromScope, strings.Join(targets, ","), cfg.ToScope, cfg.Root, cfg.DryRun)
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

func itemKindLabel(kind sync.ItemKind) string {
	switch kind {
	case sync.Instructions:
		return "instructions"
	case sync.Skills:
		return "skills"
	default:
		return "all"
	}
}

func resolveScope(specific, fallback string) (config.Scope, error) {
	s := specific
	if s == "" {
		s = fallback
	}
	if s == "" {
		return config.ScopeLocal, nil
	}
	return config.ParseScope(s)
}

func buildSyncConfig(fromStr, toStr string, dryRun bool, scope, fromScopeStr, toScopeStr string) (*config.SyncConfig, error) {
	from, err := config.ParseAgent(fromStr)
	if err != nil {
		return nil, err
	}

	fromScope, err := resolveScope(fromScopeStr, scope)
	if err != nil {
		return nil, err
	}
	toScope, err := resolveScope(toScopeStr, scope)
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
		// Only skip when same agent AND same scope
		if a == from && fromScope == toScope {
			fmt.Fprintf(os.Stderr, "warning: skipping %s (same agent and scope)\n", a)
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

	if fromScope == config.ScopeGlobal && toScope == config.ScopeGlobal && flagRoot != "" && flagRoot != "." {
		fmt.Fprintf(os.Stderr, "warning: --root is ignored when both scopes are global\n")
	}

	return &config.SyncConfig{
		From:      from,
		To:        targets,
		Root:      root,
		FromScope: fromScope,
		ToScope:   toScope,
		DryRun:    dryRun,
		Verbose:   flagVerbose,
	}, nil
}
