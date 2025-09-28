package cli

import (
	"fmt"
	"github.com/mr3iscuit/mzipext/constants"
	commands__merge_extract "github.com/mr3iscuit/mzipext/pkg/zip"

	"github.com/spf13/cobra"
)

// mergeExtractCmd represents the mergeExtract command
var mergeExtractCmd = &cobra.Command{
	Use:     constants.MergeExtractCmdUse,
	Short:   constants.MergeExtractCmdShort,
	Long:    constants.MergeExtractCmdLong,
	Example: constants.MergeExtractCmdExample,
	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {
		inputDir, _ := cmd.Flags().GetString("input-dir")
		outputDir, _ := cmd.Flags().GetString("output-dir")

		if len(args) == 0 && inputDir == "" {
			return fmt.Errorf(
				"input files are not provided: %s",
				"You should provide directory of zip folders, or provide zip path's via command line arguments.",
			)
		}

		err := commands__merge_extract.MergeExtract(
			inputDir,
			outputDir,
			cmd.Flags().Args(),
		)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeExtractCmd)
	mergeExtractCmd.Flags().String(
		"output-dir",
		"./",
		"",
	)
	mergeExtractCmd.Flags().String(
		"input-dir",
		"",
		"",
	)
}
