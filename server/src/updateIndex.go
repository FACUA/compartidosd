package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
)

func updateIndex() error {
	newIndex, err := ioutil.ReadFile(*indexPathPtr)

	if err != nil {
		return err
	}

	if json.Unmarshal(newIndex, &index); err != nil {
		return err
	}

	// If the DNS was specified, resolve the addresses
	for i := range index {
		index[i].Host = resolveHost(index[i].Host)
	}

	// Re-marshal the data again to compute the checksum
	json, err := json.Marshal(index)
	if err != nil {
		return err
	}

	sum := sha256.Sum256(json)
	indexChecksum = hex.EncodeToString(sum[:])

	return nil
}
