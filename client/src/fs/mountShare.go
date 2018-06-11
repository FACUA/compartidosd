package fs

import (
	"fmt"
	"time"

	"facua.org/compartidosd/client/util"
	"facua.org/compartidosd/common"
)

// MountShare attempts to mount a share, or returns an error
func MountShare(entry common.IndexEntry) error {
	mountPoint := getMountPoint(entry)

	err := createFolder(mountPoint)
	if err != nil {
		return err
	}

	_, err = runMountCommand(
		entry,
		sharedFolderName,
		mountPoint,
		"guest,uid=1000,iocharset=utf8,vers=1.0")

	if err != nil {
		fmt.Println("Guest mounting failed, attempting to authenticate " +
			"as guest%guest...")

		// Windows shares don't support the guest attribute,
		// but if they are public, they will allow to use any user,
		// so we use guest%guest
		_, err = runMountCommand(
			entry,
			sharedFolderName,
			mountPoint,
			"user=guest,password=guest,uid=1000,iocharset=utf8,vers=1.0")
	}

	if err != nil {
		// If the mounting fails, remove the directory, but don't care about
		// errors
		util.Cmd("rmdir", mountPoint)
	}

	return err
}

func runMountCommand(
	entry common.IndexEntry,
	sharedFolderName string,
	mountPoint string,
	options string) (string, error) {
	return util.CmdWithTimeout(
		util.CmdWithTimeoutOptions{
			Timeout: 5 * time.Second, LogOutput: true},
		"mount",
		"-t",
		"cifs",
		fmt.Sprintf("//%s/%s", entry.Host, sharedFolderName),
		mountPoint,
		"-o",
		options)
}
