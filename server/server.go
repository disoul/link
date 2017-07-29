package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LinkResponse struct {
	ErrorCode uint32
	ErrorMsg  string
	Data      interface{}
}

func handleRegister(w http.ResponseWriter, r *http.Request) LinkError {
	type registerModel struct {
		Id       string
		TypeName string
		Address  string
	}

	var body registerModel
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return LinkError{JSON_DECODE_ERROR, "Json Decoder: Can not decode request boay", err}
	}

}

// CreateLinkServer create link base http server
func CreateLinkServer(port uint32) {
	http.HandleFunc("/register", LinkHTTPHandle(handleRegister))

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatalf("Can not listen a server on %v\nError: %s", port, err)
	}
}
