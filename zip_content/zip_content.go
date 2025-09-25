package zip_content

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ZipContent struct {
	ZipPath  string
	Checksum string
	Files    []*zip.File
	r        *zip.ReadCloser
}

func NewZipContent(zipPath string) (*ZipContent, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		_ = r.Close()
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

		fmt.Printf(
			"Just a directory. Dest dir made: %s\n",
			filePath,
		)

		return err
	}

	// create destination directories
	destDir = filepath.Dir(filePath)
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

	fmt.Printf(
		"Dest dir made: %s\nFile path is %s\n",
		destDir,
		filePath,
	)

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
			"failed to copy data for file %s: %w",
			file.Name,
			err,
		)
	}

	return nil
}

func (z *ZipContent) IsSameZip(zipContent ZipContent) bool {
	return z.Checksum == zipContent.Checksum
}

func PeekZipContent(filePath string) (ZipContent, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return ZipContent{}, err
	}
	defer zipReader.Close()

	files := make(
		[]*zip.File,
		0,
	)

	for _, f := range zipReader.File {
		if f.FileInfo().Name() == "" && f.FileInfo().Size() == 0 {
			continue
		}
		files = append(
			files,
			f,
		)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return ZipContent{}, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(
		hash,
		file,
	); err != nil {
		return ZipContent{}, err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))

	return ZipContent{
		ZipPath:  filePath,
		Checksum: checksum,
		Files:    files,
	}, nil
}
