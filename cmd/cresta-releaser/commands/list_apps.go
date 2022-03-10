package commands

import (
	"github.com/spf13/cobra"
	"os"
)

var listAppsCmd = &cobra.Command{
	Use:     "apps",
	Short:   "Returns all applications",
	Example: "cresta-releaser list apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		releases, err := api.ListApplications()
		cobra.CheckErr(err)
		return getOutputFormat().WriteStringSlice(os.Stdout, releases)
	},
	Args: cobra.NoArgs,
}

func init() {
	listCmd.AddCommand(listAppsCmd)
}
