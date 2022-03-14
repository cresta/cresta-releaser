package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/cresta/cresta-releaser/releaser"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cresta-releaser",
	Short:   "Help deploy new releases",
	PreRunE: verifyOutputFormat,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg := zap.NewProductionConfig()
		cfg.Encoding = "console"
		if *verbose {
			cfg.Level.SetLevel(zap.DebugLevel)
		}
		logger, err := cfg.Build()
		if err != nil {
			return err
		}
		api, err = releaser.NewFromCommandLine(cmd.Context(), logger, nil)
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
}

var outputFormat *string
var verbose *bool

func init() {
	outputFormat = rootCmd.PersistentFlags().StringP("output", "o", "", "Output format of the command")
	verbose = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "If true, will print out verbose logging")
}
