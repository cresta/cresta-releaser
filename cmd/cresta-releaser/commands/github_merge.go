package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var githubMergeCmd = &cobra.Command{
	Use:     "merge",
	Short:   "Merge a pull request for the current repository",
	Example: "cresta-releaser github merge 1121",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("missing pull request number and/or message")
		}
		prNumber := args[0]
		prAsInt, err := strconv.ParseInt(prNumber, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid pull request number: %s", prNumber)
		}
		return api.MergePullRequestForCurrentRemote(cmd.Context(), prAsInt)
	},
	Args: cobra.ExactValidArgs(1),
}

func init() {
	githubCmd.AddCommand(githubMergeCmd)
}
