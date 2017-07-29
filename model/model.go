package model

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"link/server"
	"net/http"
	"time"
)

// ModelState model state
type ModelState uint8

const (
	STATE_IDLE ModelState = iota
	STATE_READY
	STATE_OFFLINE
	STATE_ERROR
)

// ModelType model type
type ModelType struct {
	name string
}

// Model link communication base unit
type Model struct {
	_id       string
	modelType ModelType
	address   string

	state    ModelState
	lastbeat int32
}

// NewModel Model constructor
func NewModel(typename, address, id string) (Model, error) {
	// 直接提供id说明为重启的Model
	if id != "" {
		// TODO: search redis model info
	}

	modelType := ModelType{typename}

	model := Model{
		modelType: modelType,
		address:   address,
		state:     STATE_IDLE,
	}

	// http communication test
	initData := map[string]string{
		"data": "ping",
	}

	res, err := model.Send(initData)
	if err != nil {
		return model, err
	}

	dataMap, _ := server.MapDecode(res.Body)
	// test not pass
	if dataMap["data"] != "pong" {
		var err error
		err = "Error: ping model address can not get expect response"
		return model, err
	}

	return model, nil
}

// Send send message to model
func (model Model) Send(data interface{}) (*http.Response, error) {
	byteData, err := json.Marshal(data)
	if err != nil {
		return &http.Response{}, err
	}

	res, err := http.Post(model.address, "application/json", bytes.NewReader(byteData))
	if err != nil {
		return &http.Response{}, err
	}

	return res, nil
}

func genModelID(model Model) string {
	s := fmt.Sprintf("%s%s%v", model.modelType.name, model.address, time.Now().Unix())
	h := sha1.New()
	h.Write([]byte(s))

	return fmt.Sprintf("%x".h.Sum(nil))
}
