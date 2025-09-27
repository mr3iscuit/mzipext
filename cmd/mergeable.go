package cmd

import (
	"fmt"
	commands__mergeable "github.com/mr3iscuit/mzipext/pkg/zip"

	"github.com/mr3iscuit/mzipext/constants"
	"github.com/spf13/cobra"
)

var mergeableCmd = &cobra.Command{
	Use:          constants.MergeableCmdUse,
	Short:        constants.MergeableCmdShort,
	Long:         constants.MergeableCmdLong,
	SilenceUsage: true,

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {
		zipFiles := cmd.Flags().Args()

		ok, err := commands__mergeable.Mergeable(zipFiles)
		if err != nil {
			return err
		}

		if ok {
			fmt.Printf("Files can be merged")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeableCmd)
}
