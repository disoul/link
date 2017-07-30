package server

import "net/http"
import "fmt"

// LinkErrorCode errCode
type LinkErrorCode uint16

// LinkError error wrapper with code and msg
type LinkError struct {
	ErrorCode LinkErrorCode
	ErrorMsg  string
	error     error
}

func (error LinkError) Error() string {
	return fmt.Sprintf("LinkError: %s\nError: %s", error.ErrorMsg, error.error.Error())
}

// LinkHTTPHandle handle error to custom ServerHTTP
type LinkHTTPHandle func(http.ResponseWriter, *http.Request) LinkError

// JSON_DECODE_ERROR parse request or response body error
const (
	JSON_ENCODE_ERROR LinkErrorCode = 4000
	JSON_DECODE_ERROR LinkErrorCode = 4001
	MODEL_INIT_ERROR  LinkErrorCode = 4002
	REDIS_SAVE_ERROR  LinkErrorCode = 4003
	UNEXPECT_MESSAGE  LinkErrorCode = 4004
)

func (fn LinkHTTPHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)

	if err.Error != nil {
		res, e := LinkResponse{
			ErrorCode: err.ErrorCode,
			ErrorMsg:  err.ErrorMsg,
		}.Encode()
		if e != nil {
			http.Error(w, e.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		http.Error(w, res, 500)
	}
}
