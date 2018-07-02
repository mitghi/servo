package main

import (
	"io/ioutil"
	"strings"
)

// - MARK: Function section.

func ParseURI(uri string) []string {
	return strings.Split(uri, "/")
}

func URIGetHead(uri string) (string, bool) {
	var (
		cl         int
		components []string
	)
	components = ParseURI(uri)
	cl = len(components)
	if cl == 0 || cl <= 1 {
		return __EMPTY__, false
	}
	for _, c := range components {
		if len(c) > 0 {
			return c, true
		}
	}
	return __EMPTY__, false
}

// - MARK: File Utilities section.

func readFile(path string) ([]byte, error) {
	var (
		content []byte
		err     error
	)
	content, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return content, nil
}
