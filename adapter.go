package rboot

import (
	"fmt"
)

type Adapter interface {
	Name() string           // 适配器名称
	Incoming() chan Message // 接收到的消息
	Outgoing() chan Message // 回复的消息
}

type adapterF func(*Robot) Adapter

var adapters = make(map[string]adapterF)

// RegisterAdapter 注册适配器，名称不可重复
// 适配器需实现 Adapter 接口
func RegisterAdapter(name string, adp adapterF) {
	if name == "" {
		panic("RegisterAdapter: adapter must have a name")
	}
	if _, ok := adapters[name]; ok {
		panic("RegisterAdapter: adapter named " + name + " already registered. ")
	}
	adapters[name] = adp
}

// DetectAdapter 根据适配器名称获取适配器实例
func DetectAdapter(name string) (adapterF, error) {
	if adp, ok := adapters[name]; ok {
		return adp, nil
	}

	if len(adapters) == 0 {
		return nil, fmt.Errorf("no adapter available")
	}

	if name == "" {
		if len(adapters) == 1 {
			for _, adp := range adapters {
				return adp, nil
			}
		}
		return nil, fmt.Errorf("multiple adapters available; must choose one")
	}
	return nil, fmt.Errorf("unknown adapter '%s'", name)
}
