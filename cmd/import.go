package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
	"github.com/LaneBirmingham/coding-agent-sync/internal/sync"
	"github.com/spf13/cobra"
)

func newImportCmd() *cobra.Command {
	var (
		flagTo     string
		flagScope  string
		flagInput  string
		flagDryRun bool
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import agent config from a ZIP archive",
		Long:  "Import instructions and skills from a ZIP archive to one or more agents.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doImport(flagTo, flagScope, flagInput, flagDryRun)
		},
	}

	cmd.Flags().StringVar(&flagTo, "to", "", "destination agent(s), comma-separated")
	cmd.Flags().StringVarP(&flagScope, "scope", "", "local", "scope (local, global)")
	cmd.Flags().StringVarP(&flagInput, "input", "i", "", "input ZIP path")
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "preview import without writing")

	_ = cmd.MarkFlagRequired("to")
	_ = cmd.MarkFlagRequired("input")

	return cmd
}

func doImport(toStr, scopeStr, input string, dryRun bool) error {
	scope, err := config.ParseScope(scopeStr)
	if err != nil {
		return err
	}

	var targets []config.Agent
	for _, s := range strings.Split(toStr, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		a, err := config.ParseAgent(s)
		if err != nil {
			return err
		}
		targets = append(targets, a)
	}

	if len(targets) == 0 {
		return fmt.Errorf("no valid destination agents specified")
	}

	root := flagRoot
	if root == "" || root == "." {
		root, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}
	}

	cfg := &config.ImportConfig{
		To:     targets,
		Root:   root,
		Scope:  scope,
		Input:  input,
		DryRun: dryRun,
	}

	if flagVerbose {
		targetNames := make([]string, 0, len(targets))
		for _, t := range targets {
			targetNames = append(targetNames, string(t))
		}
		fmt.Fprintf(os.Stderr, "verbose: import to=%s scope=%s root=%s input=%s dry-run=%t\n", strings.Join(targetNames, ","), cfg.Scope, cfg.Root, cfg.Input, cfg.DryRun)
	}

	result, err := sync.Import(cfg)
	if err != nil {
		return err
	}

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}

	for _, action := range result.Actions {
		fmt.Println(action)
	}

	return nil
}
