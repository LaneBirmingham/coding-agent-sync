package cmd

import (
	"github.com/LaneBirmingham/coding-agent-sync/internal/sync"
	"github.com/spf13/cobra"
)

func newDiffCmd() *cobra.Command {
	var (
		flagFrom string
		flagTo   string
	)

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Preview sync changes (alias for sync --dry-run)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doSync(sync.All, flagFrom, flagTo, true)
		},
	}

	cmd.Flags().StringVar(&flagFrom, "from", "", "source agent")
	cmd.Flags().StringVar(&flagTo, "to", "", "destination agent(s)")

	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")

	return cmd
}
