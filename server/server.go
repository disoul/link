package server

import (
	"fmt"
	"link/model"
	"link/storage"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

var linkDatabase *mgo.Database

func handleRegister(w http.ResponseWriter, r *http.Request) LinkError {
	type registerModel struct {
		Id       string `json:"id"`
		TypeName string `json:"type"`
		Address  string `json:"address"`
	}
	var data registerModel

	err := LinkDecode(r.Body, &data)
	if err != nil {
		return LinkError{JSON_DECODE_ERROR, "Json Decoder: Can not decode request boay", err}
	}

	model, err := model.NewModel(data.TypeName, data.Address, data.Id)
	if err != nil {
		return LinkError{MODEL_INIT_ERROR, "Register Model: init Model fail", err}
	}

	id, err := storage.UpsertModel(model, linkDatabase)
	if err != nil {
		return LinkError{REDIS_SAVE_ERROR, "Register Model: save to redis faild", err}
	}

	w.Header().Set("Content-Type", "application/json")
	res, err := LinkResponse{Data: id}.Encode()
	if err != nil {
		return LinkError{JSON_ENCODE_ERROR, "Json Encoder: Can not encode response boay", err}
	}
	w.Write([]byte(res))

	return LinkError{error: nil}
}

// sendMessage core handle in link
// send messages between 2 models
func sendMessage(w http.ResponseWriter, r *http.Request) LinkError {
	type LinkMessage struct {
		targetID    string
		targetQuery map[string]string
		originID    string
		message     map[string]string
	}
	var message LinkMessage

	err := LinkDecode(r.Body, &message)
	if err != nil {
		return LinkError{UNEXPECT_MESSAGE, "can not decode message with LinkMessage", err}
	}

	return LinkError{error: nil}
}

// CreateLinkServer create link base http server
func CreateLinkServer(db *mgo.Database, port uint32) {
	linkDatabase = db

	http.HandleFunc("/register", LinkHTTPHandle(handleRegister).ServeHTTP)
	http.HandleFunc("/send", LinkHTTPHandle(sendMessage).ServeHTTP)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatalf("Can not listen a server on %v\nError: %s", port, err)
	}
}
