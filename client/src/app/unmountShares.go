package app

import (
	"facua.org/compartidosd/client/fs"
	"facua.org/compartidosd/common"
)

func unmountShares(shares []share) (unmounted []share, errored []share) {
	done := make(chan bool, len(shares))

	for _, share := range shares {
		// For some reason we can't use the "share" alias here, we get
		// "share" is not a type
		go func(share common.IndexEntry) {
			err := fs.UnmountShare(share)

			if err == nil {
				unmounted = append(unmounted, share)
			} else {
				errored = append(errored, share)
			}

			done <- true
		}(share)
	}

	// Wait until all operations are finished
	for i := 0; i < len(shares); i++ {
		<-done
	}

	return unmounted, errored
}
