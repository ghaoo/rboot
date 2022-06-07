package rboot

import (
	"fmt"
)

var (
	plugins = make(map[string]Plugin)
	ruleset = make(map[string]map[string]string)
)

// Plugin 脚本结构体
type Plugin struct {
	Action      PluginFunc        // 执行解析或一些必要加载
	Ruleset     map[string]string // 脚本规则集合
	Usage       map[string]string // 帮助信息
	Description string            // 简介
}

// PluginFunc SetupFunc 脚本执行或解析函数
// - bot: A Robot instance
// - incoming: The incoming message
type PluginFunc func(bot *Robot, incoming *Message) []*Message

// RegisterPlugin 注册插件
func RegisterPlugin(name string, plugin Plugin) {

	if name == "" {
		panic("RegisterPlugin: plugin must have a name")
	}

	// 如果已经存在插件，将插件替换为现在的插件
	plugins[name] = plugin

	if len(plugin.Ruleset) > 0 {

		ruleset[name] = plugin.Ruleset
	}
}

// DirectivePlugin 根据脚本名称获取插件执行函数
func DirectivePlugin(name string) (PluginFunc, error) {

	if plugin, ok := plugins[name]; ok {
		return plugin.Action, nil
	}

	return nil, fmt.Errorf("DirectivePlugin: no action found in plugin '%s' (missing plugin?)", name)
}

// helpSetup 帮助脚本
func helpSetup(_ *Robot, in *Message) (msg []*Message) {
	rule := in.Header.Get("rule")
	args := in.Header["args"]
	msgtype := in.Header.Get("msgtype")

	switch rule {
	case `help`:
		if len(args) < 2 || args[1] == "" {
			// 获取所有插件信息
			content := ""

			for name, plug := range plugins {
				if msgtype == "markdown" {
					content += fmt.Sprintf("**%s** - %s\n", name, plug.Description)
					for _rule, _explain := range plug.Usage {
						content += fmt.Sprintf("- `%s` - %s\n\n", _rule, _explain)
					}

				} else {
					content += fmt.Sprintf("- %s - %s\n", name, plug.Description)
				}
			}
			msg = append(msg, NewMessage(content))
		} else {
			if plug, ok := plugins[args[1]]; ok {
				usage := ""
				for _rule, _explain := range plug.Usage {
					if msgtype == "markdown" {
						usage += fmt.Sprintf("> `%s` - %s \n", _rule, _explain)
					} else {
						usage += fmt.Sprintf("- %s - %s \n", _rule, _explain)
					}
				}

				msg = append(msg, NewMessage(usage))
			} else {
				msg = append(msg, NewMessage("未找到插件 "+args[1]+" \n"))
			}
		}
	}

	return msg
}

// 帮助脚本规则集
var helpRules = map[string]string{
	`help`: `^!help(?: *)(\S*)`,
}

func init() {
	RegisterPlugin(`help`, Plugin{
		Action:  helpSetup,
		Ruleset: helpRules,
		Usage: map[string]string{
			"!help":          "查看所有插件",
			"!help <plugin>": "查看plugin插件的帮助信息",
		},
		Description: `显示已经注册的所有插件或帮助信息`,
	})
}
