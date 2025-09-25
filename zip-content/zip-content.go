package zip_content

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

type ZipContent struct {
	FilePath string
	Checksum string
	Files    []zip.FileHeader
}

func (z ZipContent) IsSameFile(zipContent ZipContent) bool {
	return z.Checksum == zipContent.Checksum
}

func PeekZipContent(filePath string) (ZipContent, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return ZipContent{}, err
	}
	defer zipReader.Close()

	files := make([]zip.FileHeader, 0)

	for _, f := range zipReader.File {
		if f.FileInfo().Name() == "" && f.FileInfo().Size() == 0 {
			continue
		}
		files = append(files, f.FileHeader)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return ZipContent{}, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ZipContent{}, err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))

	return ZipContent{
		FilePath: filePath,
		Checksum: checksum,
		Files:    files,
	}, nil
}
