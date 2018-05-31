package rboot

import (
	"fmt"
)

var (
	availablePlugins = make(map[string]*Plugin)

	availableAdapters = make(map[string]func(*Controller)Adapter)
)

type Plugin struct {
	Rules       []string  // 插件指令集
	Action      Handler // 插件操作函数
	Description string    // 插件简介
}

func ListPlugins() map[string][]string {
	var ps = make(map[string][]string)

	for name, plug := range availablePlugins {
		ps[name] = plug.Rules
	}

	return ps
}

type Handler func(*Controller) ([]Message, error)

// 注册插件
func RegisterPlugin(name string, plugin *Plugin) {
	if name == "" {
		panic("RegisterPlugin: plugin must have a name")
	}
	if _, ok := availablePlugins[name]; ok {
		panic("RegisterPlugin: plugin named " + name + " already registered. ")
	}

	availablePlugins[name] = plugin
}

func DirectiveAction(name string) (Handler, error) {

	if plugin, ok := availablePlugins[name]; ok {
		return plugin.Action, nil
	}

	return nil, fmt.Errorf("DirectiveAction: no action found in plugin '%s' (missing a plugin?)", name)

}

type Adapter interface {
	Run() error // 运行适配器
	Emit(*Controller) error	// 将 in message 发送给控制器
	Send(*Controller, ...string) error	// 发送消息
	Reply(*Controller, ...string) error	// 回答
}

func RegisterConnecter(name string, f func(*Controller) Adapter) {
	availableAdapters[name] = f
}
