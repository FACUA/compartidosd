package network

import (
	"fmt"
	"time"

	"facua.org/compartidosd/client/util"
	"facua.org/compartidosd/common"
)

// OnlineStatusQuery represents the result of the QueryEntriesOnlineStatus
// function
type OnlineStatusQuery struct {
	Online  []common.IndexEntry
	Offline []common.IndexEntry
}

func isHostOnline(host string) (bool, error) {
	_, err := util.CmdWithTimeout(
		util.CmdWithTimeoutOptions{
			Timeout: 500 * time.Millisecond, LogOutput: false},
		"ping", "-c", "1", host)

	if err == nil {
		return true, nil
	}

	if err.Error() == "exit status 1" || err.Error() == "Command timed out" {
		return false, nil
	}

	return false, err
}

// QueryEntriesOnlineStatus pings all given entries in parallel and returns
// which are online and which are offline.
func QueryEntriesOnlineStatus(
	entries []common.IndexEntry) OnlineStatusQuery {
	var online []common.IndexEntry
	var offline []common.IndexEntry

	done := make(chan bool, len(entries))

	for _, entry := range entries {
		go func(entry common.IndexEntry) {
			isOnline, err := isHostOnline(entry.Host)
			if err != nil {
				// Ping failed, assume that it's offline
				// (isOnline will be false in case of error)
				fmt.Println(
					"Failed to ping " + entry.Host + ":\n" + err.Error())
			}

			if isOnline {
				online = append(online, entry)
			} else {
				offline = append(offline, entry)
			}

			done <- true
		}(entry)
	}

	// Wait until all queries are finished
	for i := 0; i < len(entries); i++ {
		<-done
	}

	return OnlineStatusQuery{Online: online, Offline: offline}
}
