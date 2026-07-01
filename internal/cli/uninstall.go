package cli

import (
	"fmt"

	"github.com/coderianx/mira/internal/install"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <repo>",
	Short: "Remove an installed tool",
	Long: `Remove a tool installed by mira.

Examples:
  mira uninstall github.com/user/reponame
  mira uninstall github.com/user/reponame@v1.0.0`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := install.Remove(args[0]); err != nil {
			return fmt.Errorf("uninstall: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
