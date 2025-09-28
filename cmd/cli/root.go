package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   RootCommandUse,
	Short: RootCommandShort,
	Long:  RootCommandLong,
	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {
		if v, _ := cmd.Flags().GetBool("version"); v {
			cmd.Println(RootCommandUse + " version " + RootCommandVersion)
			os.Exit(0)
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP(
		"version",
		"v",
		false,
		"version",
	)
}
