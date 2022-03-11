package commands

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var getCmd = &cobra.Command{
	Use:     "git",
	Short:   "Git operations for new release cuts",
	Example: "cresta-releaser get release customer-namespace 00-staging",
}

func init() {
	rootCmd.AddCommand(getCmd)
}
