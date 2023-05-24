package oxywd

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/sik0-o/0xyWD/template"
)

// SetProxy build a proxy extension for browser and than inject it in browser capabilities.
// It also creates temporary folder `tmp`.
// in this version of method empty proxy do not fires an error.
func SetProxy(caps CapsExtAdder, proxy string) error {
	if proxy == "" {
		// TODO: or ERROR?
		return nil
	}

	// createExtension
	tempLocation := fmt.Sprintf("tmp/%s/", Md5Str(proxy))
	if err := os.MkdirAll(tempLocation, 0777); err != nil {
		return err
	}

	if err := createExtension(tempLocation, proxy, true); err != nil {
		return err
	}

	zipfilename := fmt.Sprintf("tmp/%s.zip", Md5Str(proxy))
	if err := createZipFromFolder(tempLocation, zipfilename); err != nil {
		return err
	}

	return caps.AddExtension(zipfilename)
}

func createExtension(temporaryExtensionLocation string, proxy string, buildintemplates bool) error {
	p, err := url.Parse(proxy)
	if err != nil {
		return err
	}

	if buildintemplates {
		return createFromBuiltinTemplates(temporaryExtensionLocation, map[string]string{
			"${PRX_PROXY}": p.String(),
		})
	} else {
		if err := createFromTemplates(`template`, temporaryExtensionLocation, map[string]string{
			"${PRX_PROXY}": p.String(),
		}); err != nil {
			return err
		}
	}

	return nil
}

func createFromBuiltinTemplates(extensionLocation string, replaceData map[string]string) error {
	for _, filename := range []string{
		"manifest.js",
		"background.js",
	} {
		str := template.Get(filename)
		for k, v := range replaceData {
			str = strings.ReplaceAll(str, k, v)
		}

		b := []byte(str)
		if err := ioutil.WriteFile(extensionLocation+"/"+filename, b, 0777); err != nil {
			return err
		}
	}

	return nil
}

func createFromTemplates(templatesLocation string, extensionLocation string, replaceData map[string]string) error {
	return filepath.Walk(templatesLocation, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			// create directory in temporary location
		} else {
			// create file in temporary location
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			str := string(b)

			for k, v := range replaceData {
				str = strings.ReplaceAll(str, k, v)
			}

			b = []byte(str)
			if err := ioutil.WriteFile(extensionLocation+"/"+info.Name(), b, 0777); err != nil {
				return err
			}
		}

		return nil
	})
}

// createZipFromFolder создает zip архив zipfilename из директории location и возвращает в случае ошибки error
func createZipFromFolder(location string, zipfilename string) error {
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
