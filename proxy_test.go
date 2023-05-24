package oxywd

import (
	"io/ioutil"
	"os"
	"testing"
)

var cleanupTemporaryDir = false

func TestCreateZipFromFolder(t *testing.T) {
	t.Log("Start createZipFromFolder Test")
	testFolder, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("testFolder %s was created", testFolder)

	// add temp files
	if tmpFile, err := createTempFileWithContent(testFolder, "testA"); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Temporary file `%s` was created", tmpFile)
	}
	if tmpFile, err := createTempFileWithContent(testFolder, "testB"); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Temporary file `%s` was created", tmpFile)
	}
	if tmpFile, err := createTempFileWithContent(testFolder, "testC"); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Temporary file `%s` was created", tmpFile)
	}

	testZipLocationDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	testZipLocation := testZipLocationDir + string(os.PathSeparator) + "test.zip"
	t.Logf("testZipLocation %s was calculated", testFolder)

	t.Log("Start creating zip from folder", testFolder, testZipLocation)
	if err := createZipFromFolder(testFolder, testZipLocation); err != nil {
		t.Fatal(err)
	}

	if zipFileInfo, err := os.Stat(testZipLocation); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("`%s` stats{size:%d modTime:%s}", zipFileInfo.Name(), zipFileInfo.Size(), zipFileInfo.ModTime().Local().String())
	}

	if cleanupTemporaryDir {
		// remove temporary files from test
		if err := os.RemoveAll(testFolder); err != nil {
			t.Fatal(err)
		}
		if err := os.RemoveAll(testZipLocationDir); err != nil {
			t.Fatal(err)
		}
	}

	t.Logf("Nice, have a good day!")
}

func createTempFileWithContent(location string, content string) (string, error) {
	f, err := ioutil.TempFile(location, "")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return "", err
	}

	return f.Name(), nil
}
