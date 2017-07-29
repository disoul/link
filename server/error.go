package server

import "net/http"

// LinkErrorCode errCode
type LinkErrorCode uint16

// LinkError error wrapper with code and msg
type LinkError struct {
	ErrorCode LinkErrorCode
	ErrorMsg  string
	Error     error
}

// LinkHTTPHandle handle error to custom ServerHTTP
type LinkHTTPHandle func(http.ResponseWriter, *http.Request) LinkError

// JSON_DECODE_ERROR parse request or response body error
const JSON_DECODE_ERROR LinkErrorCode = 4001

func (fn LinkHTTPHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)

	if err.Error != nil {
		http.Error(w, err.Error(), 500)
	}
}
