package cli

import (
	"github.com/spf13/cobra"

	"github.com/coderianx/mira/internal/output"
)

const Version = "v1.0.0"

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		output.Box("mira "+Version, output.Cyan)
		output.Dim("  github.com/coderianx/mira")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
