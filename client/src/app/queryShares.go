package app

import (
	"facua.org/compartidosd/client/api"
)

var lastIndex []share

// queryShares is the first step of the tick. It takes care of queying the
// server using the API, and comparing the results to those of the last tick.
// It will return the shares that were added, and the shares that were deleted
// since the last tick.
// If this is the first tick, then all the shares returned by the server will
// be returned as new, and no deleted shares will be returned.
func queryShares(mountedShares []share) (
	newShares []share,
	deletedShares []share,
	err error) {
	newIndex, err := api.GetIndex()

	newShares = make([]share, 0)
	deletedShares = make([]share, 0)

	if err != nil {
		return newShares, deletedShares, err
	}

	if len(newIndex) == 0 {
		return newShares, deletedShares, nil
	}

	for _, share := range newIndex {
		if !contains(lastIndex, share) {
			newShares = append(newShares, share)
		}
	}

	for _, share := range mountedShares {
		if !contains(newIndex, share) {
			deletedShares = append(deletedShares, share)
		}
	}

	lastIndex = newIndex

	return newShares, deletedShares, nil
}

func contains(array []share, share share) bool {
	for _, shareToCompare := range array {
		if shareToCompare.Name == share.Name {
			return true
		}
	}

	return false
}
