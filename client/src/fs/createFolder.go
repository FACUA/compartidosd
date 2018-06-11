package fs

import (
	"os"

	"facua.org/compartidosd/client/util"
)

func createFolder(path string) error {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	_, err = util.Cmd("chown", "root:"+networkGroup, path)
	if err != nil {
		return err
	}

	_, err = util.Cmd("chmod", "775", path)
	if err != nil {
		return err
	}

	return nil
}
