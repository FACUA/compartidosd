package fs

import (
	"io/ioutil"
	"strings"
)

type bookmarkFile struct {
	Path     string
	Contents string
}

func getBookmarkFiles() (files []bookmarkFile, err error) {
	homeDirs, err := ioutil.ReadDir("/home")
	if err != nil {
		return files, err
	}

	for _, homeDir := range homeDirs {
		path := "/home/" + homeDir.Name() + "/.config/gtk-3.0/bookmarks"
		bytes, err := ioutil.ReadFile(path)

		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				continue
			} else {
				return files, err
			}
		}

		contents := string(bytes)
		files = append(files, bookmarkFile{path, contents})
	}

	return files, nil
}
