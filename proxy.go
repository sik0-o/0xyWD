package oxywd

import (
	"fmt"
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
	if err := СreateZipFromFolder(tempLocation, zipfilename); err != nil {
		return err
	}

	return caps.AddExtension(zipfilename)
}

func createExtension(temporaryExtensionLocation string, proxy string, buildintemplates bool) error {
	p, err := url.Parse(proxy)
	if err != nil {
		return err
	}

	proxyURL := p.String()
	if strings.Contains(proxyURL, "://") {
		if s := strings.Split(proxyURL, "://"); len(s) > 1 {
			proxyURL = s[1]
		}
	}

	data := map[string]string{
		"${PRX_PROXY}":    proxyURL,
		"${EXT_FULLNAME}": "SikO.o/0xyWD_PRX",
	}

	if buildintemplates {
		return createFromBuiltinTemplates(temporaryExtensionLocation, data)
	} else if err := createFromTemplates(`template`, temporaryExtensionLocation, data); err != nil {
		return err
	}

	return nil
}

func createFromBuiltinTemplates(extensionLocation string, replaceData map[string]string) error {
	for _, filename := range []string{
		"manifest.json",
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
