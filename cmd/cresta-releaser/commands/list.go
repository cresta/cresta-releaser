package commands

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Returns all instances of an object",
	Example: "cresta-releaser list releases",
}

func init() {
	rootCmd.AddCommand(listCmd)
}
