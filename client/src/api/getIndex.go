package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"facua.org/compartidosd/client/fs/cache"
	"facua.org/compartidosd/common"
)

var (
	serverAddress = os.Getenv("SERVER_ADDRESS")
	lastChecksum  string
)

// GetIndex returns the latest Index. If the index hasn't changed since the
// last time we queried it, it returns an index of length zero.
func GetIndex() ([]common.IndexEntry, error) {
	var index common.Index

	checksum, err := getChecksum()

	if err != nil {
		return defaultToCache(err)
	}

	if checksum == lastChecksum {
		return index.Entries, nil
	}

	resp, err := httpClient.Get(serverAddress)
	if err != nil {
		return defaultToCache(err)
	}

	if resp.StatusCode != http.StatusOK {
		return defaultToCache(
			errors.New("The status code is not HTTP 200 OK"))
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return defaultToCache(err)
	}

	if json.Unmarshal(bytes, &index); err != nil {
		return defaultToCache(err)
	}

	lastChecksum = index.Checksum
	go func() {
		err := cache.WriteIndex(index)
		if err != nil {
			fmt.Println("Failed to write the index cache:")
			fmt.Println(err)
		}
	}()

	return index.Entries, nil
}

func defaultToCache(err error) ([]common.IndexEntry, error) {
	fmt.Println("Couldn't read the index from the server:")
	fmt.Println(err)
	fmt.Println("Trying to read the index from the cache...")

	index, err := cache.ReadIndex()

	if err == nil && index.Checksum == "" {
		return index.Entries, errors.New("There is no cache saved")
	}

	return index.Entries, err
}
