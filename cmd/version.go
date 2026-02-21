package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "dev"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("cas v%s\n", Version)
		},
	}
}
