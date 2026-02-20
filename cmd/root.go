package cmd

import (
	"github.com/spf13/cobra"
)

var (
	flagRoot    string
	flagVerbose bool
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "cas",
		Short: "Sync configuration between coding agents",
		Long:  "cas (coding-agent-sync) syncs instructions and skills between Claude Code, GitHub Copilot, and OpenCode.",
	}

	root.PersistentFlags().StringVar(&flagRoot, "root", ".", "project root directory")
	root.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")

	root.AddCommand(newSyncCmd())
	root.AddCommand(newDiffCmd())
	root.AddCommand(newVersionCmd())

	return root
}
