package cli

import (
	"fmt"

	"github.com/coderianx/mira/internal/install"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <repo>",
	Short: "Install a tool from a GitHub repo",
	Long: `Install a tool by reading its mira.json manifest.

Examples:
  mira install github.com/user/reponame
  mira install github.com/user/reponame@v1.0.0`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := install.Install(args[0]); err != nil {
			return fmt.Errorf("install: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
