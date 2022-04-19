package commands

import (
	"github.com/spf13/cobra"
)

var releaseApplyCmd = &cobra.Command{
	Use:     "apply",
	Short:   "Apply a release",
	Example: "cresta-releaser release apply customer-namespace 01-prod-alpha",
	RunE: func(cmd *cobra.Command, args []string) error {
		oldRelease, newRelease, err := api.PreviewRelease(cmd.Context(), args[0], args[1], false)
		cobra.CheckErr(err)
		return api.ApplyRelease(args[0], args[1], oldRelease, newRelease)
	},
	Args: cobra.ExactValidArgs(2),
}

func init() {
	releaseCmd.AddCommand(releaseApplyCmd)
}
