package cli

import (
	"fmt"

	"github.com/coderianx/mira/internal/install"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <repo>",
	Short: "Reinstall a tool to get the latest version",
	Long: `Re-download and reinstall a tool.

Examples:
  mira update github.com/user/reponame
  mira update github.com/user/reponame@v1.0.0`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := install.Update(args[0]); err != nil {
			return fmt.Errorf("update: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
