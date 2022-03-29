package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var githubPrCmd = &cobra.Command{
	Use:     "pr",
	Short:   "Pull request of curent branch",
	Example: "cresta-releaser github pr",
	RunE: func(cmd *cobra.Command, args []string) error {
		pr, err := api.PullRequestCurrent(cmd.Context())
		cobra.CheckErr(err)
		return getOutputFormat().WriteObject(os.Stdout, pr)
	},
	Args: cobra.NoArgs,
}

func init() {
	githubCmd.AddCommand(githubPrCmd)
}
