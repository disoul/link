package server

import (
	"fmt"
	"link/model"
	"link/storage"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2"
)

var linkDatabase *mgo.Database

func handleRegister(w http.ResponseWriter, r *http.Request) LinkError {
	type registerModel struct {
		ID       int32  `json:"id"`
		TypeName string `json:"type"`
		Address  string `json:"address"`
	}
	var data registerModel

	err := LinkDecode(r.Body, &data)
	if err != nil {
		return LinkError{JSON_DECODE_ERROR, "Json Decoder: Can not decode request boay", err}
	}

	model, err := model.NewModel(data.TypeName, data.Address, data.ID)
	if err != nil {
		return LinkError{MODEL_INIT_ERROR, "Register Model: init Model fail", err}
	}

	id, err := storage.UpsertModel(model, linkDatabase)
	if err != nil {
		return LinkError{REDIS_SAVE_ERROR, "Register Model: save to redis faild", err}
	}

	res, err := LinkResponse{Data: id}.Encode()
	if err != nil {
		return LinkError{JSON_ENCODE_ERROR, "Json Encoder: Can not encode response boay", err}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, res)

	return LinkError{}
}

// sendMessage core handle in link
// send messages between 2 models
func sendMessage(w http.ResponseWriter, r *http.Request) LinkError {
	type LinkMessage struct {
		Limit       int                    `json:"limit"`
		TargetID    int32                  `json:"targetId"`
		TargetQuery map[string]interface{} `json:"targetQuery"`
		OriginID    int32                  `json:"originId"`
		Message     map[string]interface{} `json:"message"`
	}
	var message LinkMessage

	err := LinkDecode(r.Body, &message)
	if err != nil {
		return LinkError{UNEXPECT_MESSAGE, "can not decode message with LinkMessage", err}
	}

	send := func(m model.Model, c chan [2]string) {
		res, err := m.Send(message.Message)
		var result [2]string
		result[0] = strconv.Itoa(int(m.ID))
		result[1] = "ok"
		if err != nil {
			result[1] = fmt.Sprintf("http error:%s", err.Error())
		}
		if res.StatusCode != 200 {
			result[1] = fmt.Sprintf(LinkDecodeString(res.Body))
		}

		c <- result
	}

	models, err := storage.FindModels(message.TargetQuery, message.TargetID, message.Limit, linkDatabase)
	responses := make(map[string]string)
	c := make(chan [2]string)

	fmt.Println(models, len(models))

	if len(models) > 0 {
		for _, m := range models {
			go send(m, c)
		}
		chanIndex := 0
		for i := range c {
			responses[i[0]] = i[1]
			chanIndex++
			if chanIndex++; chanIndex >= len(models) {
				close(c)
			}
		}
	}

	fmt.Println(responses)

	res, err := LinkResponse{Data: responses}.Encode()
	if err != nil {
		return LinkError{JSON_ENCODE_ERROR, "can not encode response data", err}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, res)

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
