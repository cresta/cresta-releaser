package commands

import (
	"github.com/spf13/cobra"
)

var gitCommitCmd = &cobra.Command{
	Use:     "commit",
	Short:   "Create a commit message for the current release",
	Example: "cresta-releaser git commit customer-namespace 00-staging",
	RunE: func(cmd *cobra.Command, args []string) error {
		cobra.CheckErr(api.CommitForRelease(cmd.Context(), args[0], args[1]))
		return nil
	},
	Args: cobra.ExactValidArgs(2),
}

func init() {
	gitCmd.AddCommand(gitCommitCmd)
}
