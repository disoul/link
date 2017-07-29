package model

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"link/utils"
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
	Name string
}

// Model link communication base unit
type Model struct {
	ID        string
	ModelType ModelType
	Address   string

	State    ModelState
	Lastbeat int32
}

// NewModel Model constructor
func NewModel(typename, address, id string) (Model, error) {
	// 直接提供id说明为重启的Model
	modelType := ModelType{typename}

	model := Model{
		ModelType: modelType,
		Address:   address,
		State:     STATE_IDLE,
		ID:        id,
	}
	if id == "" {
		model.genModelID()
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
	if fmt.Sprintf("%x", dataMap["data"]) != "pong" {
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

// Stringify stringify model
func (model Model) Stringify() (string, error) {
	bytesData, err := json.Marshal(model)
	if err != nil {
		return "", err
	}

	return string(bytesData[:]), nil
}

func (model *Model) genModelID() {
	s := fmt.Sprintf("%s%s%v", model.ModelType.Name, model.Address, time.Now().Unix())
	h := sha1.New()
	h.Write([]byte(s))

	model.ID = string(h.Sum(nil)[:])
}
