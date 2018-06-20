package rboot

import (
	"fmt"
)

var (
	availableScripts = make(map[string]*Script)

	availableConnecters = make(map[string]func(*Response) Connecter)

	rulesets = make(map[string]map[string]string)
)

type Script struct {
	Action      SetupFunc         // 插件操作函数
	Ruleset     map[string]string // 指令集
	Description string            // 插件简介
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

type Connecter interface {
	Name() string          // 适配器名称
	Run() error            // 运行适配器
	Send(...string) error  // 发送消息
	Reply(...string) error // 回复消息
	Close() error          // 关闭适配器
}

func RegisterConnecter(name string, f func(*Response) Connecter) {
	availableConnecters[name] = f
}

func getConnecter(res *Response, name string) (Connecter, error) {
	if c, ok := availableConnecters[name]; ok {
		return c(res), nil
	}

	if len(availableConnecters) == 0 {
		return nil, fmt.Errorf("no connecter available")
	}

	if name == "" {
		if len(availableConnecters) == 1 {
			for _, c := range availableConnecters {
				return c(res), nil
			}
		}
		return nil, fmt.Errorf("multiple connecters available; must choose one")
	}
	return nil, fmt.Errorf("unknown connecter '%s'", name)
}