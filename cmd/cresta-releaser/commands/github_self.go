package commands

import (
	"github.com/spf13/cobra"
	"os"
)

var githubSelfCmd = &cobra.Command{
	Use:     "self",
	Short:   "Show who you are on GitHub",
	Example: "cresta-releaser github self",
	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := api.GithubWhoami(cmd.Context())
		cobra.CheckErr(err)
		return getOutputFormat().WriteString(os.Stdout, user)
	},
	Args: cobra.NoArgs,
}

func init() {
	githubCmd.AddCommand(githubSelfCmd)
}
