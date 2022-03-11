package commands

import (
	"github.com/spf13/cobra"
)

var githubPrCmd = &cobra.Command{
	Use:     "pr",
	Short:   "Pull request of curent branch",
	Example: "cresta-releaser github pr",
	RunE: func(cmd *cobra.Command, args []string) error {
		cobra.CheckErr(api.PullRequestCurrent(cmd.Context()))
		return nil
	},
	Args: cobra.NoArgs,
}

func init() {
	githubCmd.AddCommand(githubPrCmd)
}
