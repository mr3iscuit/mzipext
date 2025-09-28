package cli

var (
	RootCommandUse   = "mzipext"
	RootCommandShort = "tool description"
	RootCommandLong  = "long tool desc"

	MergeableCmdUse   = "mergeable"
	MergeableCmdShort = "check if zip files are mergeable"
	MergeableCmdLong  = `
	gdextract mergeable <zip1> <zip2>

	checks if zip1 and zip2 are mergeable.
	zip1 and zip2 must be valid zip files.
	`

	MergeExtractCmdUse     = "merge-extract"
	MergeExtractCmdShort   = "extract zip files and merge into directory"
	MergeExtractCmdLong    = ""
	MergeExtractCmdExample = "mzipext --input-dir \"My Zips Folder\" --output-dir \"My Folder\" this.zip that.zip"
)
