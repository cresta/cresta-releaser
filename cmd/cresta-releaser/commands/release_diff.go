package commands

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"os"
)

var releaseDiffCmd = &cobra.Command{
	Use:     "diff",
	Short:   "Diff what would change in a release",
	Example: "cresta-releaser release diff customer-namespace 00-staging",
	RunE: func(cmd *cobra.Command, args []string) error {
		oldRelease, newRelease, err := api.PreviewRelease(args[0], args[1])
		cobra.CheckErr(err)
		oldContent, newContent := oldRelease.Yaml(), newRelease.Yaml()
		d := diffmatchpatch.New()
		diffs := d.DiffMain(oldContent, newContent, true)
		return getOutputFormat().WriteString(os.Stdout, d.DiffPrettyText(diffs))
	},
	Args: cobra.ExactValidArgs(2),
}

func init() {
	releaseCmd.AddCommand(releaseDiffCmd)
}
