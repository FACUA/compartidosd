package fs

import (
	"fmt"

	"facua.org/compartidosd/client/util"
)

// Cleanup cleans everything left by the filesystem control
func Cleanup() {
	_, err := util.Cmd("rmdir", networkFolder)
	if err != nil {
		fmt.Println("Failed to remove the network folder!")
		fmt.Println(err)
	}

	removeBookmarks()
}
