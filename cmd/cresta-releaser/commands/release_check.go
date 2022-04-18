package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var releaseCheckCmd = &cobra.Command{
	Use:     "check",
	Short:   "Check if a release would change",
	Example: "cresta-releaser release check customer-namespace 00-staging",
	RunE: func(cmd *cobra.Command, args []string) error {
		oldRelease, newRelease, err := api.PreviewRelease(cmd.Context(), args[0], args[1], true)
		cobra.CheckErr(err)
		oldContent, newContent := oldRelease.Yaml(), newRelease.Yaml()
		if oldContent != newContent {
			os.Exit(1)
		}
		os.Exit(0)
		return nil
	},
	Args: cobra.ExactValidArgs(2),
}

func init() {
	releaseCmd.AddCommand(releaseCheckCmd)
}
