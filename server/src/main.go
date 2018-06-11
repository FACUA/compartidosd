package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"facua.org/compartidosd/common"
)

var (
	index         []common.IndexEntry
	indexChecksum string
	indexPathPtr  = flag.String(
		"index",
		"./index.json",
		"The location of the index file")
	dnsServerPtr = flag.String(
		"dns",
		"",
		"The DNS server to use to resolve host names, if any")
	portPtr = flag.String(
		"port",
		"3000",
		"The TCP port to listen for HTTP connections")
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	response := common.Index{Entries: index, Checksum: indexChecksum}
	message, err := json.Marshal(response)

	if err != nil {
		panic(err)
	}

	w.Write([]byte(message))
}

func getChecksum(w http.ResponseWriter, r *http.Request) {
	response := common.ChecksumResponse{Checksum: indexChecksum}
	message, err := json.Marshal(response)

	if err != nil {
		panic(err)
	}

	w.Write([]byte(message))
}

func main() {
	flag.Parse()

	initializeIndex()

	http.HandleFunc("/", getIndex)
	http.HandleFunc("/checksum", getChecksum)
	if err := http.ListenAndServe(":"+*portPtr, nil); err != nil {
		panic(err)
	}
}
