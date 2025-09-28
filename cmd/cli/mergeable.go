package cli

import (
	"github.com/mr3iscuit/mzipext/pkg/zip"
	"github.com/spf13/cobra"
)

var mergeableCmd = &cobra.Command{
	Use:          MergeableCmdUse,
	Short:        MergeableCmdShort,
	Long:         MergeableCmdLong,
	SilenceUsage: true,

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {
		zipFiles := cmd.Flags().Args()

		inputDir, _ := cmd.Flags().GetString("input-dir")
		outputDir, _ := cmd.Flags().GetString("output-dir")

		ok, err := zip.Mergeable(
			zipFiles,
			inputDir,
			outputDir,
		)
		if err != nil {
			return err
		}

		if ok {
			cmd.Println("Files can be merged")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeableCmd)
	mergeableCmd.Flags().String(
		"output-dir",
		"./",
		"",
	)
	mergeableCmd.Flags().String(
		"input-dir",
		"",
		"",
	)
}
