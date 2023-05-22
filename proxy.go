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

	return caps.AddUnpackedExtension(tempLocation)
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

	for filename, str := range map[string]string{
		"manifest.js":   template.ManifestJS,
		"background.js": template.BackgroundJS,
	} {
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
