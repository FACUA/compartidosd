package cache

import (
	"encoding/json"
	"io/ioutil"

	"facua.org/compartidosd/common"
)

// WriteIndex stores the passed index to the index cache, and returns the error
// should one occur.
func WriteIndex(index common.Index) error {
	bytes, err := json.Marshal(index)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(indexCachePath, bytes, 0600)
	return err
}
