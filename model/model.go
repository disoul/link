package model

type ModelState unit8

const (
	STATE_IDLE ModelState = iota
	STATE_READY
	STATE_OFFLINE
	STATE_ERROR
)

type ModelType struct {
	name string
}

type Model struct {
	_id string
	_type ModelType
	address string `通讯地址`

	state ModelState     
	lastbeat int32 
}

func newModel(typename, address string): *Model {
}

