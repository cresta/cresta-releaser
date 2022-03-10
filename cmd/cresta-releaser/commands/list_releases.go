package commands

import (
	"github.com/spf13/cobra"
	"os"
)

var listReleasesCmd = &cobra.Command{
	Use:     "releases",
	Aliases: []string{"rel"},
	Short:   "Returns all release names for an application",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		apps, err := api.ListApplications()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return apps, cobra.ShellCompDirectiveDefault
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return cobra.ExactValidArgs(1)(cmd, args)
	},
	Example: "cresta-releaser list releases argo",
	RunE: func(cmd *cobra.Command, args []string) error {
		releases, err := api.ListReleases(args[0])
		cobra.CheckErr(err)
		return getOutputFormat().WriteStringSlice(os.Stdout, releases)
	},
}

func init() {
	listCmd.AddCommand(listReleasesCmd)
}
