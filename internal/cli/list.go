package cli

import (
	"github.com/coderianx/mira/internal/install"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		return install.List()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
