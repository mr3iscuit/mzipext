package constants

var (
	RootCommandUse   = "gdextract"
	RootCommandShort = "tool description"
	RootCommandLong  = "long tool desc"

	MergeableCmdUse   = "mergeable"
	MergeableCmdShort = "check if zip files are mergeable"
	MergeableCmdLong  = `
	gdextract mergeable <zip1> <zip2>

	checks if zip1 and zip2 are mergeable.
	zip1 and zip2 must be valid zip files.
	`
)
