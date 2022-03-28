package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var getReleaseCmd = &cobra.Command{
	Use:     "release",
	Short:   "Gets a release",
	Example: "cresta-releaser get release customer-namespace 00-staging",
	RunE: func(cmd *cobra.Command, args []string) error {
		releases, err := api.GetRelease(args[0], args[1])
		cobra.CheckErr(err)
		return getOutputFormat().WriteObject(os.Stdout, releases)
	},
	Args: cobra.ExactValidArgs(2),
}

func init() {
	getCmd.AddCommand(getReleaseCmd)
}
