package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:       "list",
	Aliases:   []string{"ls"},
	Short:     "Returns all instances of an object",
	Example:   "list releases",
	ValidArgs: []string{"releases"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
