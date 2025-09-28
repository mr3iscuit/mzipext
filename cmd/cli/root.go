package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   RootCommandUse,
	Short: RootCommandShort,
	Long:  RootCommandLong,
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
