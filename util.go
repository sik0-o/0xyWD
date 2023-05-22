package oxywd

import (
	"crypto/md5"
	"fmt"
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
