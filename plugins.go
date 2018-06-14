package rboot

import (
	"fmt"
)

var (
	availablePlugins = make(map[string]*Plugin)

	availableConnecters = make(map[string]func(*Response) Connecter)

	rulesets = make(map[string]map[string]string)
)

type Plugin struct {
	Action      SetupFunc // 插件操作函数
	Ruleset     map[string]string  // 指令集
	Description string    // 插件简介
}

type SetupFunc func(*Response) error

// 注册插件
func RegisterPlugin(name string, plugin *Plugin) {

	if name == "" {
		panic("RegisterPlugin: plugin must have a name")
	}
	if _, ok := availablePlugins[name]; ok {
		panic("RegisterPlugin: plugin named " + name + " already registered. ")
	}

	availablePlugins[name] = plugin

	if len(plugin.Ruleset) > 0 {

		rulesets[name] = plugin.Ruleset
	}
}

func getPlugin(name string) (*Plugin, error) {
	if plug, ok := availablePlugins[name]; ok {
		return plug, nil
	}

	if len(availablePlugins) == 0 {
		return nil, fmt.Errorf("no plug-ins available")
	}

	if name == "" {
		if len(availablePlugins) == 1 {
			for _, plug := range availablePlugins {
				return plug, nil
			}
		}
		return nil, fmt.Errorf("multiple plugins available; must choose one")
	}
	return nil, fmt.Errorf("unknown plugin '%s'", name)

}

func DirectiveAction(name string) (SetupFunc, error) {

	if plugin, ok := availablePlugins[name]; ok {
		return plugin.Action, nil
	}

	return nil, fmt.Errorf("DirectiveAction: no action found in plugin '%s' (missing a plugin?)", name)

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
