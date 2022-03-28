package commands

import (
	"github.com/spf13/cobra"
)

var gitFreshCmd = &cobra.Command{
	Use:     "fresh",
	Short:   "Fresh branch for new operations",
	Example: "cresta-releaser git fresh customer-namespace 00-staging",
	RunE: func(cmd *cobra.Command, args []string) error {
		cobra.CheckErr(api.FreshGitBranch(cmd.Context(), args[0], args[1], *forcedName))
		return nil
	},
	Args: cobra.ExactValidArgs(2),
}

var forcedName *string

func init() {
	gitCmd.AddCommand(gitFreshCmd)
	forcedName = rootCmd.PersistentFlags().StringP("name", "n", "", "Forced name for this new branch")
}
