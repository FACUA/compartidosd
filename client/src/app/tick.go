package app

import (
	"fmt"

	"facua.org/compartidosd/client/network"
	"facua.org/compartidosd/common"
)

type share = common.IndexEntry

var (
	// The shares that are already mounted on our filesystem
	mountedShares []share
	// The shares that we know that are online, but couldn't mount successfully
	// for some reason
	erroredShares []share
	// The shares that we have knowledge that they exist, but are offline
	offlineShares []share
)

// Tick performs an app tick
// A tick is an action that is performed periodically.
//
// A tick takes care of:
//    * Querying the available network shares from the server
//    * Checking their online status
//    * Mounting new shares if:
//        * They were added to the server index
//        * They were on the index, but were offline before, and came online
//          on this tick.
//    * Unmounting existing shares if:
//        * They were removed from the index
//        * They were previously online and disconnected
//
// Any of these steps might error. In that case, the next tick will re-attempt
// the step again.
func Tick() {
	fmt.Printf("=== BEGINNING TICK ===\n\n")

	// Request the share list from the server if it has changed
	fmt.Printf("== Querying shares...\n\n")
	newShares, deletedShares, err := queryShares(mountedShares)
	reportQuerySharesResult(newShares, deletedShares, err)

	// From the new, offline, and errored shares, determine...
	fmt.Printf("== Querying shares online status...\n\n")
	query := network.QueryEntriesOnlineStatus(
		append(newShares, append(offlineShares, erroredShares...)...))
	// ...the ones we have to mount...
	sharesToMount := query.Online
	// ...and the ones we'll track on the next tick
	offlineShares = query.Offline
	// Clear the erroredShares, as they are distributed between sharesToMount
	// and offlineShares
	erroredShares = []share{}

	// From the already mounted shares, determine which ones we have to unmount
	sharesToUnmount := network.QueryEntriesOnlineStatus(mountedShares).Offline
	sharesToUnmount = append(sharesToUnmount, deletedShares...)
	reportQueryOnlineStatusResult(
		sharesToMount, sharesToUnmount, offlineShares)

	// Mount the online new shares and the existing shares that came online
	fmt.Printf("== Mounting shares...\n\n")
	mounted, errored := mountShares(sharesToMount)
	mountedShares = append(mountedShares, mounted...)
	erroredShares = append(erroredShares, errored...)
	reportMountSharesResult(mounted, errored)

	// Unmount the online shares that disconnected, and the shares that were
	// deleted on the server
	fmt.Printf("== Unmounting shares...\n\n")
	unmounted, errored := unmountShares(sharesToUnmount)
	for _, share := range unmounted {
		deleteShareFromArray(&mountedShares, share)
	}
	reportUnmountSharesResult(unmounted, errored)

	fmt.Printf("\n=== TICK FINISHED  ===\n")
}

func reportQuerySharesResult(
	newShares []share,
	deletedShares []share,
	err error) {
	if err != nil {
		fmt.Println("An error occurred while querying the shares:")
		fmt.Println(err)
		fmt.Println("Continuing assuming that no shares were added or deleted.")
	}

	logShares("New shares", newShares)
	logShares("Deleted shares", deletedShares)
}

func reportQueryOnlineStatusResult(
	sharesToMount []share,
	sharesToUnmount []share,
	offlineShares []share) {

	logShares("Shares to mount", sharesToMount)
	logShares("Shares to unmount", sharesToUnmount)
	logShares("Offline shares", offlineShares)
}

func reportMountSharesResult(
	mounted []share,
	errored []share) {
	logShares("Successfully mounted shares", mounted)
	logShares("Errored shares", errored)
}

func reportUnmountSharesResult(
	unmounted []share,
	errored []share) {
	logShares("Successfully unmounted shares", unmounted)
	logShares("Errored shares", errored)
}

func logShares(groupName string, shares []share) {
	fmt.Println(groupName + ":")
	for _, share := range shares {
		fmt.Println(share.Name)
	}
	fmt.Println()
}

func deleteShareFromArray(shares *[]share, share share) {
	// Find the index of the share in the array
	i := -1

	for j := range *shares {
		if (*shares)[j].Name == share.Name {
			i = j
			break
		}
	}

	if i == -1 {
		// The share is not in the array; do nothing
		return
	}

	// Replace the item to delete with the last item
	(*shares)[i] = (*shares)[len(*shares)-1]
	// Delete the last item
	*shares = (*shares)[:len(*shares)-1]
}
