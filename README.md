# mzipext

A simple CLI tool to extract multiple ZIP files and merge them into a single directory.

## Why mzipext?

When you download large files from the internet, they often come split into multiple ZIP parts (e.g., `Archive-1.zip`, `Archive-2.zip`, ...).  
Manually extracting and merging these parts is tedious:

- Each part usually extracts into its own subfolder
- You end up with duplicated directory structures
- Verifying whether two ZIP files can be merged requires manual inspection

`mzipext` automates this process:

- **One command** to extract multiple ZIPs and merge them into a clean directory
- **Mergeability check** to ensure ZIPs has no conflicting files before extraction
- **Consistent output** so you don’t have to fix folder structures by hand

In short: it saves time, avoids errors, and simplifies working with multi-part ZIP archives.


```bash
Usage:
  mzipext [command]

Available Commands:
  completion    Generate the autocompletion script for the specified shell
  help          Help about any command
  merge-extract extract zip files and merge into directory
  mergeable     check if zip files are mergeable

Flags:
  -h, --help   help for mzipext

Use "mzipext [command] --help" for more information about a command.
```

## Available Commands

- **completion**  
  Generate the autocompletion script for the specified shell.

- **help**  
  Show help information about any command.

- **merge-extract**  
  Extract ZIP files and merge them into a target directory.  
  Useful for handling multi-part ZIP archives.

- **mergeable**  
  Check if two ZIP files are mergeable.  
  Both files must be valid ZIP archives.


## Features

- Extract ZIP files into a target folder
- Merge multiple ZIP parts into one directory
- Check if two ZIP files are mergeable

## Installation

```bash
go install github.com/mr3iscuit/mzipext@latest
```

## Usage 
```bash
mzipext merge-extract --input-dir "My Zips Folder" --output-dir "My Folder" this.zip that.zip
```

Options
--input-dir directory containing zip files
--output-dir target directory (default: ./)

## Check mergability
```bash
mzipext mergeable zip1.zip zip2.zip
```

### Example 
```bash
mzipext merge-extract --output-dir ~/Res Resources-20250918T161444Z-1-002.zip
```