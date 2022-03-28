package commands

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var releaseCmd = &cobra.Command{
	Use:     "release",
	Short:   "Operate on a release",
	Example: "cresta-releaser release check customer-namespace 00-staging",
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}
