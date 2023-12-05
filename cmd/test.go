package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

// CmdTest defines the CLI sub-command 'test'.
var CmdTest = &cobra.Command{
	Use:   "test [flags]",
	Short: "...",
	Args:  cobra.NoArgs,
	Run:   runTest,
}

// Initialize CLI options.
func init() {
}

// runTest is called when the "test" command is used.
func runTest(_ *cobra.Command, _ []string) {
	slog.Info("Done")
}
