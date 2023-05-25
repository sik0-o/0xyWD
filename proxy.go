package oxywd

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/sik0-o/0xyWD/template"
)

// setupProxyExtension creates temporary folder in `tmp/` folder than build templates
// and add it to temporary folder.
func setupProxyExtension(proxy string) (string, error) {
	if proxy == "" {
		// TODO: or ERROR?
		return "", nil
	}

	// createExtension
	tempLocation := fmt.Sprintf("tmp/%s/", Md5Str(proxy))
	if err := os.MkdirAll(tempLocation, 0777); err != nil {
		return "", err
	}

	if err := createExtension(tempLocation, proxy, true); err != nil {
		return "", err
	}

	return tempLocation, nil
}

// SetProxy build a proxy unpacked extension for browser
// and than inject it in browser capabilities.
// It also creates temporary folder `tmp`.
// in this version of method empty proxy do not fires an error.
func SetProxy(caps CapsExtAdder, proxy string) error {
	tempLocation, err := setupProxyExtension(proxy)
	if err != nil {
		return err
	}

	return caps.AddUnpackedExtension(tempLocation)
}

// SetProxyZip build a proxy ZIP extension for browser
// and than inject it in browser capabilities.
// It also creates temporary folder `tmp`.
// in this version of method empty proxy do not fires an error.
func SetProxyZip(caps CapsExtAdder, proxy string) error {
	tempLocation, err := setupProxyExtension(proxy)
	if err != nil {
		return err
	}

	if tempLocation == "" {
		return errors.New("empty tempLocation path")
	}

	zipfilename := fmt.Sprintf("tmp/%s.zip", Md5Str(proxy))
	if err := СreateZipFromFolder(tempLocation, zipfilename); err != nil {
		return err
	}

	return caps.AddExtension(zipfilename)
}

// createExtension in temporaryExtensionLocation with provided
// proxy url string (scheme:// required) from builtintemplates if set or from templates
// that are located in `template` folder.
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

	data := map[string]any{
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

func createFromBuiltinTemplates(extensionLocation string, replaceData map[string]any) error {
	for _, filename := range []string{
		"manifest.json",
		"background.js",
	} {
		ext := template.Get(filename)
		if err := ioutil.WriteFile(extensionLocation+"/"+filename, []byte(ext.BuildTemplate(replaceData)), 0777); err != nil {
			return err
		}
	}

	return nil
}

func createFromTemplates(templatesLocation string, extensionLocation string, replaceData map[string]any) error {
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
				str = strings.ReplaceAll(str, k, fmt.Sprintf("%s", v))
			}

			b = []byte(str)
			if err := ioutil.WriteFile(extensionLocation+"/"+info.Name(), b, 0777); err != nil {
				return err
			}
		}

		return nil
	})
}
