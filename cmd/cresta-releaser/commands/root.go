package commands

import (
	"context"
	"fmt"
	"github.com/cresta/cresta-releaser/releaser"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cresta-releaser",
	Short:   "Help deploy new releases",
	PreRunE: verifyOutputFormat,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		api, err = releaser.NewFromCommandLine(cmd.Context(), nil)
		return err
	},
}

func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

var api releaser.Api

func verifyOutputFormat(_ *cobra.Command, _ []string) error {
	of := getOutputFormat()
	if of == nil {
		return fmt.Errorf("output format not supported")
	}
	return nil
}

func getOutputFormat() outputFormatter {
	if outputFormat == nil {
		return nil
	}
	switch *outputFormat {
	case "":
		return &NewlineFormatter{}
	case "auto":
		return &NewlineFormatter{}
	case "json":
		return &JSONFormatter{}
	default:
		panic("Invalid formatter")
	}
	return nil
}

var outputFormat *string

func init() {
	outputFormat = rootCmd.PersistentFlags().StringP("output", "o", "", "Output format of the command")
	//api = releaser.NewFromCommandLine()
}
