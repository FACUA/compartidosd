package app

import (
	"fmt"

	"facua.org/compartidosd/client/fs"
)

// Stop stops the application daemon, unmounting every share and cleaning
// everything up.
func Stop() {
	fmt.Printf("\nStopping...\n")

	for _, share := range mountedShares {
		err := fs.UnmountShare(share)

		if err != nil {
			fmt.Println("Error when unmounting share " + share.Name)
			fmt.Println(err)
		}
	}

	fs.Cleanup()
}
