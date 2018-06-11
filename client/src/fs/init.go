package fs

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	networkFolder    string
	networkGroup     string
	sharedFolderName string
)

// Init initializes the filesystem control
func Init() {
	networkFolder = os.Getenv("NETWORK_FOLDER")

	if networkFolder == "" {
		networkFolder = os.Getenv("HOME") + "/Network"
	}

	networkGroup = os.Getenv("NETWORK_GROUP")

	if networkGroup == "" {
		networkGroup = os.Getenv("USER")
	}

	sharedFolderName = os.Getenv("SHARED_FOLDER_NAME")

	if sharedFolderName == "" {
		sharedFolderName = "Shared"
	}

	if _, err := os.Stat(networkFolder); err == nil {
		// The network folder already exists. This might mean that:
		//     - The system was shut down forcefully, and we have empty folders
		//     - The daemon crashed, and we have mounted folders.
		//
		// So for each folder, we will attempt a simple "rmdir". If the
		// directories are not empty, then we will attempt "umount" then
		// "rmdir". If all succeed, then we will remove the network folder so
		// we can have a clean start. If any fail, the daemon will crash.
		subDirs, err := ioutil.ReadDir(networkFolder)

		if err != nil {
			fmt.Println("Failed to read the network folder!")
			panic(err)
		}

		for _, subDir := range subDirs {
			dir := networkFolder + "/" + subDir.Name()
			fmt.Println("Removing " + dir)

			err = os.Remove(dir)

			if err != nil {
				fmt.Println("Failed:")
				fmt.Println(err.Error())

				fmt.Println("Unmounting and removing " + dir)
				err = unmountAndDeleteDir(dir)

				if err != nil {
					fmt.Println("Failed:")
					panic(err)
				}
			}
		}

		fmt.Println("Removing " + networkFolder)
		err = os.Remove(networkFolder)

		if err != nil {
			fmt.Println("Failed to remove the network folder!")
			panic(err)
		}
	}

	err := createFolder(networkFolder)
	if err != nil {
		panic(err)
	}

	setupBookmarks()
}
