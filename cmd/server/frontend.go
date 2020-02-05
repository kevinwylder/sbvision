package main

import (
	"net/http"
	"strings"
)

// code from https://medium.com/@hau12a1/golang-http-serve-static-files-correctly-5feb98ae9da1

// FileSystem custom file system handler
type FileSystem struct {
	http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
