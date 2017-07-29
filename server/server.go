package server

import (
	"fmt"
	"link/model"
	"link/storage"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

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

	err = storage.UpdateModel(model, redisClient)
	if err != nil {
		return LinkError{REDIS_SAVE_ERROR, "Register Model: save to redis faild", err}
	}

	return LinkError{error: nil}
}

// CreateLinkServer create link base http server
func CreateLinkServer(client *redis.Client, port uint32) {
	redisClient = client

	http.HandleFunc("/register", LinkHTTPHandle(handleRegister).ServeHTTP)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatalf("Can not listen a server on %v\nError: %s", port, err)
	}
}
