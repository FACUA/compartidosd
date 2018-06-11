package fs

import (
	"os"
	"strings"
	"time"

	"facua.org/compartidosd/client/util"
	"facua.org/compartidosd/common"
)

func unmountAndDeleteDir(mountPoint string) error {
	// The share might have been busy, so retry up to five times
	var err error
	for i := 0; i < 5; i++ {
		out, err := util.CmdWithTimeout(
			util.CmdWithTimeoutOptions{
				// The timeout is 15 seconds because if a host is down, the
				// unmount operation can take that long. Don't ask me why.
				Timeout:   15 * time.Second,
				LogOutput: true},
			"umount", mountPoint)
		if err == nil {
			break
		}

		// If we're being terminated because of an OS shutdown, the OS itself
		// might have unmounted the shares for us, so we discard the error.
		if strings.Contains(out, "not mounted") {
			err = nil
			break
		}

		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return err
	}

	err = os.Remove(mountPoint)

	return err
}

// UnmountShare attempts to unmount and delete share, retrying if necessary,
// or returns an error
func UnmountShare(entry common.IndexEntry) error {
	return unmountAndDeleteDir(getMountPoint(entry))
}
