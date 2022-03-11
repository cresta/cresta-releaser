package commands

import (
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git operations for new release cuts",
}

func init() {
	rootCmd.AddCommand(gitCmd)
}
