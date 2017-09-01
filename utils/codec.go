package utils

import (
	"encoding/json"
	"io"
)

// MapDecode io.Reader to Map
func MapDecode(body io.Reader) (map[string]json.RawMessage, error) {
	var objmap map[string]json.RawMessage
	err := json.NewDecoder(body).Decode(&objmap)
	if err != nil {
		return objmap, err
	}

	return objmap, nil
}
