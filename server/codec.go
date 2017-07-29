package server

import (
	"encoding/json"
	"io"
)

// LinkResponse response
type LinkResponse struct {
	ErrorCode LinkErrorCode `json:"errorCode"`
	ErrorMsg  string        `json:"errorMsg"`
	Data      interface{}   `json:"data"`
}

// Encode encode res to json
func (res LinkResponse) Encode() (string, error) {
	bytes, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// LinkDecode io.Reader to LinkResponse
func LinkDecode(body io.Reader) (LinkResponse, error) {
	var res LinkResponse
	err := json.NewDecoder(body).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
