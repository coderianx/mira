package cli

import (
	"fmt"

	"github.com/coderianx/mira/internal/install"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <repo>",
	Short: "Show details about an installed tool",
	Long: `Show installation details and manifest info.

Examples:
  mira info github.com/user/reponame
  mira info github.com/user/reponame@v1.0.0`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := install.Info(args[0]); err != nil {
			return fmt.Errorf("info: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
