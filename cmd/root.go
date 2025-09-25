package cmd

import (
	"os"

	"github.com/mr3iscuit/mzipext/constants"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   constants.RootCommandUse,
	Short: constants.RootCommandShort,
	Long:  constants.RootCommandLong,
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	return nil
	// },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
