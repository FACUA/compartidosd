package fs

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func removeBookmarks() {
	bookmarks, err := getBookmarkFiles()

	for _, bookmarkFile := range bookmarks {
		bookmark := "file://" + networkFolder

		if !strings.Contains(bookmarkFile.Contents, bookmark) {
			// The bookmarks file doesn't contain the Network bookmark
			continue
		}

		// Remove the bookmark from the file
		pieces := strings.Split(bookmarkFile.Contents, "\n")
		newPieces := make([]string, len(pieces)-1)

		for _, piece := range pieces {
			if piece != bookmark {
				newPieces = append(newPieces, piece)
			}
		}

		newContents := strings.Join(newPieces, "\n")
		err = ioutil.WriteFile(bookmarkFile.Path, []byte(newContents), 0664)

		if err != nil {
			fmt.Println("Cannot write to " + bookmarkFile.Path)
			fmt.Println(err)
			continue
		}
	}
}
