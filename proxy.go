package oxywd

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
)

func SetProxy(caps CapsExtAdder, proxy string) error {
	if proxy == "" {
		// TODO: or ERROR?
		return nil
	}

	// createExtension
	tempLocation := fmt.Sprintf("tmp/%s/", Md5Str(proxy))

	if err := createExtension(tempLocation, proxy); err != nil {
		return err
	}

	return caps.AddUnpackedExtension(tempLocation)
}

func createExtension(temporaryExtensionLocation string, proxy string) error {
	p, err := url.Parse(proxy)
	if err != nil {
		return err
	}

	if err := createFromTemplates(`template`, temporaryExtensionLocation, map[string]string{
		"${PRX_PROXY}": p.String(),
	}); err != nil {
		return err
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
