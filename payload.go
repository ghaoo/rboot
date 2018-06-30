package rboot

import (
	"fmt"
)

var (
	scripts = make(map[string]*Script)

	providers = make(map[string]Provider)

	memorizers = make(map[string]Memorizer)

	execCall = make(map[string]SetupFunc)

	rulesets = make(map[string]map[string]string)
)

type Script struct {
	Action      SetupFunc // 执行解析或一些必要加载
	Ruleset     map[string]string // 指令集
	Call        SetupFunc      // 直接调用运行
	Usage       string    // 使用方法
	Description string    // 简介
}

type SetupFunc func(*Robot) error

// 注册脚本
func RegisterScript(name string, script *Script) {

	if name == "" {
		panic("RegisterScript: script must have a name")
	}
	if _, ok := scripts[name]; ok {
		panic("RegisterScript: script named " + name + " already registered. ")
	}

	scripts[name] = script

	if len(script.Ruleset) > 0 {

		rulesets[name] = script.Ruleset
	}

	if script.Call != nil {
		execCall[name] = script.Call
	}
}

func DirectiveAction(name string) (SetupFunc, error) {

	if script, ok := scripts[name]; ok {
		return script.Action, nil
	}

	return nil, fmt.Errorf("DirectiveAction: no action found in script '%s' (missing a script?)", name)

}

type Provider interface {
	Name() string
	Incoming() chan Message
	Send(...Message) error   // 发送消息
	Reply(...Message) error  // 回复消息
	Close() error           // 关闭适配器
}

// 注册消息适配器
func RegisterProvider(name string, prov func() Provider) {
	if name == "" {
		panic("RegisterProvider: provider must have a name")
	}
	if _, ok := providers[name]; ok {
		panic("RegisterProvider: provider named " + name + " already registered. ")
	}
	providers[name] = prov()
}

// get provider by name
func DetectProv(name string) (Provider, error) {
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

