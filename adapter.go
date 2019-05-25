package rboot

import (
	"fmt"
)

type Adapter interface {
	Name() string
	Incoming() chan Message // 接收到的消息
	Send(Message) error      // 回复的消息

}

var adapters = make(map[string]func(*Rboot) Adapter)

func RegisterAdapter(name string, adapter func(*Rboot) Adapter) {
	if name == "" {
		panic("RegisterAdapter: adapter must have a name")
	}
	if _, ok := adapters[name]; ok {
		panic("RegisterAdapter: adapter named " + name + " already registered. ")
	}
	adapters[name] = adapter
}

func DetectAdapter(name string) (func(robot *Rboot) Adapter, error) {
	if adapter, ok := adapters[name]; ok {
		return adapter, nil
	}

	if len(adapters) == 0 {
		return nil, fmt.Errorf("no adapter available")
	}

	if name == "" {
		if len(adapters) == 1 {
			for _, adapter := range adapters {
				return adapter, nil
			}
		}
		return nil, fmt.Errorf("multiple adapters available; must choose one")
	}
	return nil, fmt.Errorf("unknown adapter '%s'", name)
}
