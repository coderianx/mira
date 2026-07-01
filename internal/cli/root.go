package cli

import (
	"github.com/spf13/cobra"

	"github.com/coderianx/mira/internal/output"
)

var rootCmd = &cobra.Command{
	Use:           "mira",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.Fatal(err)
	}
}
