package commands__merge_extract

import (
	"fmt"
	commands__mergeable "github.com/mr3iscuit/mzipext/commands/mergeable"
	"github.com/mr3iscuit/mzipext/zip_content"
	"os"
	"path/filepath"
	"strings"
)

func MergeExtract(
	inputDir string,
	outputDir string,
	files []string,
) error {

	// join sys argv files and inputDir files
	if inputDir != "" {
		entries, err := os.ReadDir(inputDir)
		if err != nil {
			return fmt.Errorf(
				"failed to read directory: %w",
				err,
			)
		}
		for _, entry := range entries {
			// skip directories
			if entry.IsDir() {
				continue
			}

			fileName := entry.Name()

			// check for the ".zip" extension (case-insensitive and robust)
			if strings.EqualFold(
				filepath.Ext(fileName),
				".zip",
			) {
				// add to files
				files = append(
					files,
					fileName,
				)
			}
		}
	}

	// check if zip files are mergeable
	_, err := commands__mergeable.Mergeable(files)
	if err != nil {
		return err
	}

	var zips []*zip_content.ZipContent
	for _, file := range files {
		zipContent, err := zip_content.NewZipContent(file)
		if err != nil {
			return fmt.Errorf(
				"could not extract zip files %w",
				err,
			)
		}
		defer zipContent.CloseZip()
		zips = append(
			zips,
			zipContent,
		)
	}

	for _, zip := range zips {
		for _, file := range zip.Files {
			err = zip_content.ExtractFile(
				file,
				outputDir,
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
