package cmd

import (
	"github.com/LaneBirmingham/coding-agent-sync/internal/sync"
	"github.com/spf13/cobra"
)

func newDiffCmd() *cobra.Command {
	var (
		flagFrom      string
		flagTo        string
		flagScope     string
		flagFromScope string
		flagToScope   string
	)

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Preview sync changes (alias for sync --dry-run)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doSync(sync.All, flagFrom, flagTo, true, flagScope, flagFromScope, flagToScope)
		},
	}

	cmd.Flags().StringVar(&flagFrom, "from", "", "source agent")
	cmd.Flags().StringVar(&flagTo, "to", "", "destination agent(s)")
	cmd.Flags().StringVar(&flagScope, "scope", "", "set both from and to scope (local, global)")
	cmd.Flags().StringVar(&flagFromScope, "from-scope", "", "source scope (overrides --scope)")
	cmd.Flags().StringVar(&flagToScope, "to-scope", "", "destination scope (overrides --scope)")

	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")

	return cmd
}
