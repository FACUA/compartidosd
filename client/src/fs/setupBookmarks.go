package fs

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func setupBookmarks() {
	bookmarks, err := getBookmarkFiles()

	for _, bookmarkFile := range bookmarks {
		bookmark := "file://" + networkFolder

		if strings.Contains(bookmarkFile.Contents, bookmark) {
			// The bookmarks file already contains the Network bookmark
			continue
		}

		newContent := bookmark + "\n" + bookmarkFile.Contents + "\n"
		err = ioutil.WriteFile(bookmarkFile.Path, []byte(newContent), 0664)

		if err != nil {
			fmt.Println("Cannot write to " + bookmarkFile.Path)
			fmt.Println(err)
			continue
		}
	}
}
