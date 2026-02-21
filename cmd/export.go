package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/LaneBirmingham/coding-agent-sync/internal/config"
	"github.com/LaneBirmingham/coding-agent-sync/internal/sync"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	var (
		flagFrom   string
		flagScope  string
		flagOutput string
		flagDryRun bool
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export agent config to a ZIP archive",
		Long:  "Export instructions and skills from an agent to a portable ZIP archive.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doExport(flagFrom, flagScope, flagOutput, flagDryRun)
		},
	}

	cmd.Flags().StringVar(&flagFrom, "from", "", "source agent (claude, copilot, codex, opencode)")
	cmd.Flags().StringVarP(&flagScope, "scope", "", "local", "scope (local, global)")
	cmd.Flags().StringVarP(&flagOutput, "output", "o", "", "output ZIP path (auto-generated if omitted)")
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "preview export without writing")

	_ = cmd.MarkFlagRequired("from")

	return cmd
}

func doExport(fromStr, scopeStr, output string, dryRun bool) error {
	from, err := config.ParseAgent(fromStr)
	if err != nil {
		return err
	}

	scope, err := config.ParseScope(scopeStr)
	if err != nil {
		return err
	}

	root := flagRoot
	if root == "" || root == "." {
		root, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}
	}

	if output == "" {
		output = fmt.Sprintf("%s-%s-%s.zip", from, scope, time.Now().Format("20060102T150405"))
	}

	cfg := &config.ExportConfig{
		From:       from,
		Root:       root,
		Scope:      scope,
		Output:     output,
		DryRun:     dryRun,
		CASVersion: Version,
	}

	if flagVerbose {
		fmt.Fprintf(os.Stderr, "verbose: export from=%s scope=%s root=%s output=%s dry-run=%t\n", cfg.From, cfg.Scope, cfg.Root, cfg.Output, cfg.DryRun)
	}

	result, err := sync.Export(cfg)
	if err != nil {
		return err
	}

	for _, action := range result.Actions {
		fmt.Println(action)
	}

	if !dryRun {
		fmt.Fprintf(os.Stderr, "archive written to %s\n", output)
	}

	return nil
}
