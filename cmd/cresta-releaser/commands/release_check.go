package commands

import (
	"os"

	"github.com/cresta/cresta-releaser/releaser"

	"github.com/spf13/cobra"
)

var releaseCheckCmd = &cobra.Command{
	Use:     "check",
	Short:   "Check if a release would change",
	Example: "cresta-releaser release check customer-namespace 00-staging",
	RunE: func(cmd *cobra.Command, args []string) error {
		hasChanges, err := releaser.NeedsPromotion(cmd.Context(), api, args[0], args[1])
		cobra.CheckErr(err)
		if hasChanges {
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
