package commands__mergeable

import (
	"archive/zip"
	"fmt"

	"github.com/mr3iscuit/mzipext/zip_content"
)

type ErrRepeatedZipFile struct {
	Message          string
	RepeatedZipFiles []string
}

func (e *ErrRepeatedZipFile) Error() string {
	return fmt.Sprintf(
		"repeated zip files found: %v",
		e.RepeatedZipFiles,
	)
}

type ErrZipFileConfict struct {
	Message         string
	ConflictDetails map[string][]ZipAndZipFile
}

func (e *ErrZipFileConfict) Error() string {
	detailMsg := ""
	for name, conflicts := range e.ConflictDetails {
		detailMsg += fmt.Sprintf(
			"\n\tFile: %s found in multiple Zip files.",
			name,
		)

		for _, conflict := range conflicts {
			detailMsg += fmt.Sprintf(
				"\n\t\tFound in : %s, Size: %d bytes",
				conflict.ZipName,
				conflict.File.UncompressedSize64,
			)
		}
	}

	return fmt.Sprintf(
		"%s. Conflicts found: %d.\nDetails: %s",
		e.Message,
		len(e.ConflictDetails),
		detailMsg,
	)
}

func Mergeable(files []string) (bool, error) {
	var hasRepeatedZip bool

	repeated_count := make(map[string]int)

	for _, name := range files {
		count, ok := repeated_count[name]
		if !ok {
			count = 0
		}

		if count >= 1 {
			hasRepeatedZip = true
		}

		repeated_count[name] = count + 1
	}

	var repeated []string

	for name, value := range repeated_count {
		if value > 1 {
			repeated = append(
				repeated,
				name,
			)
		}
	}

	if hasRepeatedZip {
		return false, &ErrRepeatedZipFile{
			"Zip files can not be merged. Repeated Zip Files Error",
			repeated,
		}
	}

	zipContents := make(
		[]*zip_content.ZipContent,
		0,
	)

	for _, path := range files {
		zipContent, err := zip_content.NewZipContent(path)
		defer zipContent.CloseZip()
		if err != nil {
			return false, err
		}

		zipContents = append(
			zipContents,
			zipContent,
		)
	}

	isConflicting, conflictingFiles := hasFileConflict(zipContents)
	if isConflicting {
		return false, &ErrZipFileConfict{
			"Zip files has conflicting files",
			conflictingFiles,
		}
	}

	return true, nil
}

type ZipAndZipFile struct {
	ZipName string
	File    *zip.File
}

func hasFileConflict(zips []*zip_content.ZipContent) (bool, map[string][]ZipAndZipFile) {
	fileSeen := make(
		map[string][]ZipAndZipFile,
		0,
	)
	conflicting := make(
		map[string][]ZipAndZipFile,
		0,
	)
	isConflicting := false

	for _, contents := range zips {
		for _, file := range contents.Files {

			val, ok := fileSeen[file.Name]
			if !ok {
				val = make(
					[]ZipAndZipFile,
					0,
				)
			}

			val = append(
				val,
				ZipAndZipFile{contents.ZipPath, file},
			)

			fileSeen[file.Name] = val
		}
	}

	for name, v := range fileSeen {
		if len(v) > 1 {
			conflicting[name] = v
			isConflicting = true
		}
	}

	return isConflicting, conflicting
}
