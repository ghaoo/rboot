package rboot

import (
	"sync"
	"plugin"
)

type Rboot struct {
	plugs map[string]Plugin
}

type Instance struct {
	serverType string
	wg *sync.WaitGroup

	Storage   map[interface{}]interface{}
	StorageMu sync.RWMutex
}

func executeDirectives(inst *Instance) {
	//
}
