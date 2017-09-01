package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"link/utils"
	"net/http"
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
	Name string
}

// Model link communication base unit
type Model struct {
	ID        int32
	ModelType ModelType
	Address   string

	State    ModelState
	Lastbeat int32
}

// NewModel Model constructor
func NewModel(typename, address string, id int32) (Model, error) {
	// 直接提供id说明为重启的Model
	modelType := ModelType{typename}

	model := Model{
		ModelType: modelType,
		Address:   address,
		State:     STATE_IDLE,
		ID:        id,
	}
	// http communication test
	initData := map[string]string{
		"data": "ping",
	}

	res, err := model.Send(initData)
	if err != nil {
		return model, err
	}

	dataMap, _ := utils.MapDecode(res.Body)
	// test not pass
	if string(dataMap["data"]) != `"pong"` {
		return model, errors.New("Error: ping model address can not get expect response")
	}

	return model, nil
}

// Send send message to model
func (model Model) Send(data interface{}) (*http.Response, error) {
	byteData, err := json.Marshal(data)
	if err != nil {
		return &http.Response{}, err
	}

	res, err := http.Post(model.Address, "application/json", bytes.NewReader(byteData))
	if err != nil {
		return &http.Response{}, err
	}

	return res, nil
}

// Mapify mapify model to save
func (model Model) Mapify() map[string]interface{} {
	fields := make(map[string]interface{})
	fields["id"] = model.ID
	fields["type"] = model.ModelType.Name
	fields["address"] = model.Address
	fields["state"] = string(model.State)
	fields["lastbeat"] = string(model.Lastbeat)

	return fields
}
