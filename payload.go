package rboot

import (
	"fmt"
)

var (
	availableScripts = make(map[string]*Script)

	availableProviders = make(map[string]func(*Response) Provider)

	rulesets = make(map[string]map[string]string)
)

type Script struct {
	Action      SetupFunc         // 操作函数
	Ruleset     map[string]string // 指令集
	Description string            // 简介
}

type SetupFunc func(*Response) error

// 注册插件
func RegisterScript(name string, script *Script) {

	if name == "" {
		panic("RegisterScript: script must have a name")
	}
	if _, ok := availableScripts[name]; ok {
		panic("RegisterScript: script named " + name + " already registered. ")
	}

	availableScripts[name] = script

	if len(script.Ruleset) > 0 {

		rulesets[name] = script.Ruleset
	}
}

func DirectiveAction(name string) (SetupFunc, error) {

	if script, ok := availableScripts[name]; ok {
		return script.Action, nil
	}

	return nil, fmt.Errorf("DirectiveAction: no action found in script '%s' (missing a script?)", name)

}

type Provider interface {
	Name() string          // 适配器名称
	Run() error            // 运行适配器
	Send(...string) error  // 发送消息
	Reply(...string) error // 回复消息
	Close() error          // 关闭适配器
}

func RegisterProvider(name string, f func(*Response) Provider) {
	availableProviders[name] = f
}

func getProvider(res *Response, name string) (Provider, error) {
	if c, ok := availableProviders[name]; ok {
		return c(res), nil
	}

	if len(availableProviders) == 0 {
		return nil, fmt.Errorf("no connecter available")
	}

	if name == "" {
		if len(availableProviders) == 1 {
			for _, c := range availableProviders {
				return c(res), nil
			}
		}
		return nil, fmt.Errorf("multiple connecters available; must choose one")
	}
	return nil, fmt.Errorf("unknown connecter '%s'", name)
}