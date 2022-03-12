package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var githubApproveCmd = &cobra.Command{
	Use:     "approve",
	Short:   "Approve a pull request for the current repository",
	Example: "cresta-releaser github approve 1121 'This pull request is ready to be merged'",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("missing pull request number and/or message")
		}
		prNumber := args[0]
		message := args[1]
		prAsInt, err := strconv.ParseInt(prNumber, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid pull request number: %s", prNumber)
		}
		return api.ApprovePullRequestForCurrentRemote(cmd.Context(), message, prAsInt)
	},
	Args: cobra.ExactValidArgs(2),
}

func init() {
	githubCmd.AddCommand(githubApproveCmd)
}
