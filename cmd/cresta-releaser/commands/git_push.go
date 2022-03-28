package commands

import (
	"github.com/spf13/cobra"
)

var gitPushCmd = &cobra.Command{
	Use:     "push",
	Short:   "Push the current branch to the remote repository",
	Example: "cresta-releaser git push",
	RunE: func(cmd *cobra.Command, args []string) error {
		cobra.CheckErr(api.ForcePushCurrentBranch(cmd.Context()))
		return nil
	},
	Args: cobra.NoArgs,
}

func init() {
	gitCmd.AddCommand(gitPushCmd)
}
