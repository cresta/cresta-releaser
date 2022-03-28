package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var githubCheckPrCmd = &cobra.Command{
	Use:     "prcheck",
	Short:   "Check if there is already a PR for the current branch.  Prints the PR number if there is one.  Exit 1 if there is not.  Exit 2 on other errors",
	Example: "cresta-releaser github pr",
	RunE: func(cmd *cobra.Command, args []string) error {
		prNum, err := api.CheckForPROnCurrentBranch(cmd.Context())
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(2)
		}
		if prNum != 0 {
			fmt.Println(prNum)
		} else {
			os.Exit(1)
		}
		return nil
	},
	Args: cobra.NoArgs,
}

func init() {
	githubCmd.AddCommand(githubCheckPrCmd)
}
