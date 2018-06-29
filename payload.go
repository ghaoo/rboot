package rboot

import (
	"fmt"
)

var (
	scripts = make(map[string]*Script)

	providers = make(map[string]Provider)

	memorizers = make(map[string]Memorizer)

)

type Script struct {
	Action      SetupFunc   // 执行解析或一些必要加载
	Hook        func(Robot) // 钩子
	Usage       string      // 使用方法
	Description string      // 简介
}

type SetupFunc func(Robot, Message) []Message

// 注册插件
func RegisterScript(name string, script *Script) {

	if name == "" {
		panic("RegisterScript: script must have a name")
	}
	if _, ok := scripts[name]; ok {
		panic("RegisterScript: script named " + name + " already registered. ")
	}

	scripts[name] = script
}

type Provider interface {
	Incoming() chan Message
	Outgoing() chan Message
	Error() error
}

// register provider
func RegisterProvider(name string, f func()Provider) {
	if name == "" {
		panic("RegisterProvider: provider must have a name")
	}
	if _, ok := providers[name]; ok {
		panic("RegisterProvider: provider named " + name + " already registered. ")
	}
	providers[name] = f()
}

// get provider by name
func Detect(name string) (Provider, error) {
	if prov, ok := providers[name]; ok {
		return prov, nil
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no provider available")
	}

	if name == "" {
		if len(providers) == 1 {
			for _, prov := range providers {
				return prov, prov.Error()
			}
		}
		return nil, fmt.Errorf("multiple providers available; must choose one")
	}
	return nil, fmt.Errorf("unknown provider '%s'", name)
}

type Memorizer interface {
	Save(key string, value []byte)
	Read(key string) ([]byte, bool)
	Update(key string, value []byte)
	Delete(key string)
	Error() error
}

func RegisterMemorizer(name string, m Memorizer) {
	if name == "" {
		panic("RegisterMemorizer: memorizer must have a name")
	}
	if _, ok := memorizers[name]; ok {
		panic("RegisterMemorizer: memorizers named " + name + " already registered. ")
	}
	memorizers[name] = m
}
