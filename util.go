package oxywd

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type CapsExtensionAdder interface {
	AddExtension(path string) error
}

type CapsUnpackedExtensionAdder interface {
	AddUnpackedExtension(basePath string) error
}

type CapsExtAdder interface {
	CapsExtensionAdder
	CapsUnpackedExtensionAdder
}

func AddExtension(caps CapsExtensionAdder, path string) error {
	return caps.AddExtension(path)
}

func AddUnpackedExtension(caps CapsUnpackedExtensionAdder, basePath string) error {
	return caps.AddUnpackedExtension(basePath)
}

func Md5Str(hashing string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(hashing)))
}

// СreateZipFromFolder создает zip архив zipfilename из директории location
// и в случае ошибки возвращает error.
func СreateZipFromFolder(location string, zipfilename string) error {
	targetFile, err := os.Create(zipfilename)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	targetZipWriter := zip.NewWriter(targetFile)
	defer targetZipWriter.Close()

	return filepath.Walk(location, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.Base(path)

		targetItem, err := targetZipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(targetItem, file)
		if err != nil {
			return err
		}

		return nil
	})
}
