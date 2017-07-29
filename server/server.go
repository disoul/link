package server

import (
	"fmt"
	"log"
	"net/http"
)

func handleRegister(w http.ResponseWriter, r *http.Request) LinkError {
	type registerModel struct {
		Id       string
		TypeName string
		Address  string
	}

	_, err := LinkDecode(r.Body)
	if err != nil {
		return LinkError{JSON_DECODE_ERROR, "Json Decoder: Can not decode request boay", err}
	}

	return LinkError{Error: nil}
}

// CreateLinkServer create link base http server
func CreateLinkServer(port uint32) {
	http.HandleFunc("/register", LinkHTTPHandle(handleRegister).ServeHTTP)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatalf("Can not listen a server on %v\nError: %s", port, err)
	}
}
