package fs

import (
	"fmt"

	"facua.org/compartidosd/common"
)

func getMountPoint(entry common.IndexEntry) string {
	return fmt.Sprintf("%s/%s", networkFolder, entry.Name)
}
