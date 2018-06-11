package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"facua.org/compartidosd/common"
)

func getChecksum() (string, error) {
	var response common.ChecksumResponse

	resp, err := httpClient.Get(serverAddress + "/checksum")
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("The status code is not HTTP 200 OK")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if json.Unmarshal(bytes, &response); err != nil {
		return "", err
	}

	return response.Checksum, nil
}
