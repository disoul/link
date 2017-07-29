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

// LinkDecode io.Reader to json
func LinkDecode(body io.Reader, data interface{}) error {
	err := json.NewDecoder(body).Decode(data)
	if err != nil {
		return err
	}

	return nil
}

// MapDecode io.Reader to Map
func MapDecode(body io.Reader) (*map[string]*json.RawMessage, error) {
	var objmap map[string]*json.RawMessage
	err := json.NewDecoder(body).Decode(&objmap)
	if err != nil {
		return &objmap, err
	}

	return &objmap, nil
}
