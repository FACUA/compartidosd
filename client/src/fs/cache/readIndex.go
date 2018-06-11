package cache

import (
	"encoding/json"
	"io/ioutil"

	"facua.org/compartidosd/common"
)

const indexCachePath = "/etc/compartidosd/index-cache.json"

// ReadIndex reads the index from the cache and returs it, if it exists,
// otherwise returns a zero Index.
func ReadIndex() (index common.Index, err error) {
	bytes, err := ioutil.ReadFile(indexCachePath)
	if err != nil {
		if err.Error() == "no such file or directory" {
			return index, nil
		}
		return index, err
	}

	if len(bytes) == 0 {
		return index, nil
	}

	err = json.Unmarshal(bytes, &index)
	return index, err
}
