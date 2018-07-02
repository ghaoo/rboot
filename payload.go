package rboot

import (
	"fmt"
)

var (
	scripts = make(map[string]*Script)

	providers = make(map[string]func() Provider)

	memorizers = make(map[string]Memorizer)

	execCall = make(map[string]Call)
)

type Script struct {
	Action      SetupFunc // 执行解析或一些必要加载
	Call        Call      // 直接调用运行
	Usage       string    // 使用方法
	Description string    // 简介
}

type SetupFunc func(Robot, Message) []Message

// 用于直接调用运行的函数
type Call func(Robot) error

// 注册脚本
func RegisterScript(name string, script *Script) {

	if name == "" {
		panic("RegisterScript: script must have a name")
	}
	if _, ok := scripts[name]; ok {
		panic("RegisterScript: script named " + name + " already registered. ")
	}

	scripts[name] = script

	if script.Call != nil {
		execCall[name] = script.Call
	}
}

type Provider interface {
	Incoming() chan Message
	Outgoing() chan Message
	Error() error
}

// 注册消息适配器
func RegisterProvider(name string, prov func() Provider) {
	if name == "" {
		panic("RegisterProvider: provider must have a name")
	}
	if _, ok := providers[name]; ok {
		panic("RegisterProvider: provider named " + name + " already registered. ")
	}
	providers[name] = prov
}

// get provider by name
func DetectProv(name string) (func() Provider, error) {
	if prov, ok := providers[name]; ok {
		return prov, nil
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no provider available")
	}

	if name == "" {
		if len(providers) == 1 {
			for _, prov := range providers {
				return prov, nil
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

// 注册存储器
func RegisterMemorizer(name string, m func() Memorizer) {
	if name == "" {
		panic("RegisterMemorizer: memorizer must have a name")
	}
	if _, ok := memorizers[name]; ok {
		panic("RegisterMemorizer: memorizers named " + name + " already registered. ")
	}
	memorizers[name] = m()
}

// get memorizer by name
func DetectMemo(name string) (Memorizer, error) {
	if memo, ok := memorizers[name]; ok {
		return memo, memo.Error()
	}

	if len(memorizers) == 0 {
		return nil, fmt.Errorf("no memorizer available")
	}

	if name == "" {
		if len(memorizers) == 1 {
			for _, memo := range memorizers {
				return memo, nil
			}
		}
		return nil, fmt.Errorf("multiple memorizers available; must choose one")
	}
	return nil, fmt.Errorf("unknown memorizers '%s'", name)
}

