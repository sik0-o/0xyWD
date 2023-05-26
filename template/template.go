package template

import (
	"fmt"
	"strings"
)

var tmpls = (ExtensionList{
	NewManifestJSON(),
	NewBackgroundJS(),
})

type ExtensionList []Extension

func (el *ExtensionList) ToMap() ExtensionMap {
	em := map[string]Extension{}
	for _, e := range *el {
		em[e.Filename()] = e
	}

	return em
}

type ExtensionMap map[string]Extension

func (em *ExtensionMap) ToList() ExtensionList {
	el := []Extension{}
	for _, e := range *em {
		el = append(el, e)
	}

	return el
}

func Get(filename string) Extension {
	if tpl, ok := tmpls.ToMap()[filename]; ok {
		return tpl
	}

	return Extension{}
}

func GetTmpl(filename string) string {
	ext := Get(filename)
	return strings.Trim(ext.template, "\n\r\t ")
}

func GetBytes(filename string) []byte {
	return []byte(GetTmpl(filename))
}

type Extension struct {
	filename string
	template string
	dataVals map[string]any
}

func (e *Extension) BuildTemplate(values map[string]any) string {
	if e.dataVals == nil {
		e.dataVals = map[string]any{}
	}
	vals := e.dataVals
	for k, v := range values {
		vals[k] = v
	}

	str := e.template
	for k, v := range vals {
		str = strings.ReplaceAll(str, k, fmt.Sprintf("%s", v))
	}

	return str
}

func (e *Extension) Filename() string {
	return e.filename
}
