package commands

import (
	"github.com/spf13/cobra"
)

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Github operations for new release cuts",
}

func init() {
	rootCmd.AddCommand(githubCmd)
}
