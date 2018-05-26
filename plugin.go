package rboot

import "fmt"

var plugins = make(map[string]Plugin)

type Plugin struct {
	ServerType string    // 插件类型
	Action     SetupFunc // 插件操作函数
}

type SetupFunc func(ctx *Context) error

// 注册插件
func RegisterPlugin(name string, plugin Plugin) {
	if name == "" {
		panic("plugin must have a name")
	}
	if _, ok := plugins[name]; ok {
		panic("plugin named " + name + " already registered for server type " + plugin.ServerType)
	}
	plugins[name] = plugin
}

func DirectiveAction(name string) (SetupFunc, error) {

	if plugin, ok := plugins[name]; ok {
		return plugin.Action, nil
	}

	return nil, fmt.Errorf("no action found in plugin '%s' (missing a plugin?)", name)

}
