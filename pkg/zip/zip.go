package zip

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

type ErrZipFileConflict struct {
	Message         string
	ConflictDetails map[string][]SourceAndFile
}

func (e *ErrZipFileConflict) Error() string {
	detailMsg := ""
	for name, conflicts := range e.ConflictDetails {
		detailMsg += fmt.Sprintf(
			"\nFile: \"%s\"",
			name,
		)

		for _, conflict := range conflicts {
			detailMsg += fmt.Sprintf(
				"\n\tFound in: \"%s\", File size: %d bytes",
				conflict.SourceName,
				conflict.File.UncompressedSize64,
			)
		}
	}

	return fmt.Sprintf(
		"%s. found %d Conflict: %s",
		e.Message,
		len(e.ConflictDetails),
		detailMsg,
	)
}

func MergeExtract(
	inputDir string,
	outputDir string,
	files []string,
) error {
	// check if zip files are mergeable
	_, err := Mergeable(
		files,
		inputDir,
		outputDir,
	)
	if err != nil {
		return err
	}

	var zips []*ZipContent
	for _, file := range files {
		zipContent, err := NewZipContent(file)
		if err != nil {
			return fmt.Errorf(
				"could not extract zip files %w",
				err,
			)
		}
		zips = append(
			zips,
			zipContent,
		)
		_ = zipContent.CloseZip()
	}

	for _, oneZip := range zips {
		for _, file := range oneZip.Files {
			err = ExtractFile(
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

func Mergeable(
	files []string,
	inputDir string,
	outputDir string,
) (bool, error) {
	// join sys argv files and inputDir files

	var inputDirZips []fs.DirEntry
	if inputDir != "" {
		var err error
		inputDirZips, err = os.ReadDir(inputDir)
		if err != nil {
			return false, fmt.Errorf(
				"failed to read directory: %w",
				err,
			)
		}
	}

	for _, entry := range inputDirZips {
		// skip directories
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		// check for the ".zip" extension (case-insensitive and robust)
		if !strings.EqualFold(
			filepath.Ext(fileName),
			".zip",
		) {
			//todo: not a zip file, return an error
			continue
		}

		// add to files
		files = append(
			files,
			fileName,
		)
	}

	var hasRepeatedZip bool

	repeatedCount := make(map[string]int)

	for _, name := range files {
		count, ok := repeatedCount[name]
		if !ok {
			count = 0
		}

		if count >= 1 {
			hasRepeatedZip = true
		}

		repeatedCount[name] = count + 1
	}

	var repeated []string

	for name, value := range repeatedCount {
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
		[]*ZipContent,
		0,
	)

	for _, path := range files {
		zipContent, err := NewZipContent(path)
		defer zipContent.CloseZip()
		if err != nil {
			return false, err
		}

		zipContents = append(
			zipContents,
			zipContent,
		)
	}

	isConflicting, conflictingFiles, err := findConflictingFiles(
		zipContents,
		outputDir,
	)
	if err != nil {
		return false, err
	}

	if isConflicting {
		return false, &ErrZipFileConflict{
			"Zip files has conflicting files",
			conflictingFiles,
		}
	}

	return true, nil
}

type SourceFile struct {
	Name               string
	UncompressedSize64 uint64
}

type SourceAndFile struct {
	SourceName string
	File       SourceFile
}

func findConflictingFiles(
	zips []*ZipContent,
	outputDir string,
) (bool, map[string][]SourceAndFile, error) {
	fileSeen := make(
		map[string][]SourceAndFile,
	)
	conflicting := make(
		map[string][]SourceAndFile,
	)
	isConflicting := false

	if !strings.HasSuffix(
		outputDir,
		string(os.PathSeparator),
	) {
		outputDir += string(os.PathSeparator)
	}

	err := filepath.WalkDir(
		outputDir,
		func(
			path string,
			d fs.DirEntry,
			err error,
		) error {
			if err != nil {
				return err
			}

			if d.Name() == "." || d.Name() == ".." {
				return nil
			}

			if d.IsDir() {
				return nil
			}

			path = filepath.Clean(path) + string(filepath.Separator)
			filePath := filepath.Clean(path)

			rel := strings.TrimPrefix(
				filePath,
				outputDir,
			)

			fInfo, _ := d.Info()

			fileSeen[rel] = append(
				fileSeen[rel],
				SourceAndFile{
					outputDir,
					SourceFile{Name: rel, UncompressedSize64: uint64(fInfo.Size())},
				},
			)

			return nil
		},
	)
	if err != nil {
		return false, make(map[string][]SourceAndFile), fmt.Errorf(
			"error walking the directory tree: %w",
			err,
		)
	}

	for _, contents := range zips {
		for _, file := range contents.Files {

			val, ok := fileSeen[file.Name]
			if !ok {
				val = make(
					[]SourceAndFile,
					0,
				)
			}

			val = append(
				val,
				SourceAndFile{
					contents.ZipPath,
					SourceFile{Name: file.Name, UncompressedSize64: file.UncompressedSize64},
				},
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

	return isConflicting, conflicting, nil
}

type ZipContent struct {
	ZipPath  string
	Checksum string
	Files    []*zip.File
	r        *zip.ReadCloser
}

func NewZipContent(zipPath string) (*ZipContent, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return &ZipContent{}, err
	}

	files := make(
		[]*zip.File,
		0,
	)

	for _, f := range r.File {
		if f.FileInfo().Name() == "" && f.FileInfo().Size() == 0 {
			continue
		}
		files = append(
			files,
			f,
		)
	}

	file, err := os.Open(zipPath)
	if err != nil {
		return &ZipContent{}, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(
		hash,
		file,
	); err != nil {
		return &ZipContent{}, err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))

	return &ZipContent{
		ZipPath:  zipPath,
		Checksum: checksum,
		Files:    files,
		r:        r,
	}, nil
}

func (z *ZipContent) CloseZip() error {
	return z.r.Close()
}

func ExtractFile(
	file *zip.File,
	destDir string,
) error {
	// if it's a file, open the zip file's content reader
	rc, err := file.Open()
	defer rc.Close()
	if err != nil {
		return fmt.Errorf(
			"failed to open zip entry %s: %w",
			file.Name,
			err,
		)
	}

	filePath := filepath.Join(
		destDir,
		file.Name,
	)

	if file.FileInfo().IsDir() {
		err = os.MkdirAll(
			filepath.Dir(filePath),
			0755,
		)

		return err
	}

	// create destination directories
	_ = filepath.Dir(filePath)
	err = os.MkdirAll(
		filepath.Dir(filePath),
		0755,
	)
	if err != nil {
		return fmt.Errorf(
			"could not create destination directory %w",
			err,
		)
	}

	// create the destination file
	outFile, err := os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		file.Mode(),
	)
	defer outFile.Close()
	if err != nil {
		rc.Close()
		return fmt.Errorf(
			"failed to create destination file %s: %w",
			filePath,
			err,
		)
	}

	// copy the file data
	_, err = io.Copy(
		outFile,
		rc,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to copy data of file %s: %w",
			file.Name,
			err,
		)
	}

	return nil
}

func (z *ZipContent) IsSameZip(zipContent ZipContent) bool {
	return z.Checksum == zipContent.Checksum
}
