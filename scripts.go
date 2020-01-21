package rboot

import (
	"fmt"
)

var (
	scripts = make(map[string]Script)
	ruleset = make(map[string]map[string]string)
)

// Script 脚本结构体
type Script struct {
	Action      ScriptFunc        // 执行解析或一些必要加载
	Ruleset     map[string]string // 脚本规则集合
	Usage       map[string]string // 帮助信息
	Description string            // 简介
}

// SetupFunc 脚本执行或解析函数
// - bot: A Robot instance
// - incoming: The incoming message
type ScriptFunc func(bot *Robot, incoming *Message) []*Message

// RegisterScripts 注册脚本
func RegisterScripts(name string, script Script) {

	if name == "" {
		panic("RegisterScripts: the script must have a name")
	}
	if _, ok := scripts[name]; ok {
		log.Warnf("RegisterScripts: script named %s already registered, old script will be replaced", name)
	}

	scripts[name] = script

	if len(script.Ruleset) > 0 {

		ruleset[name] = script.Ruleset
	}
}

// DirectiveScript 根据脚本名称获取脚本执行函数
func DirectiveScript(name string) (ScriptFunc, error) {

	if script, ok := scripts[name]; ok {
		return script.Action, nil
	}

	return nil, fmt.Errorf("DirectiveScript: no action found in script '%s' (missing a script?)", name)
}

// helpSetup 帮助脚本
func helpSetup(bot *Robot, in *Message) (msg []*Message) {
	rule := in.Header.Get("rule")
	args := in.Header["args"]

	switch rule {
	case `help`:
		if len(args) < 2 || args[1] == "" {
			msg = append(msg, NewMessage(script()))
		} else {
			if scr, ok := scripts[args[1]]; ok {
				msgtype := in.Header.Get("msgtype")
				usage := ""
				for _rule, _explain := range scr.Usage {
					if msgtype == "markdown" {
						usage += fmt.Sprintf("> `%s` - %s \n", _rule, _explain)
					} else {
						usage += fmt.Sprintf("- %s - %s \n", _rule, _explain)
					}
				}

				msg = append(msg, NewMessage(usage))
			} else {
				msg = append(msg, NewMessage("未找到脚本 "+args[1]+" \n"))
			}
		}
	}

	return msg
}

func script() string {
	// 获取所有脚本信息
	content := ""

	for scr, spt := range scripts {
		content += fmt.Sprintf("- %s - %s\n", scr, spt.Description)
	}

	return content
}

// 帮助脚本规则集
var helpRules = map[string]string{
	`help`: `^!help(?: *)(\S*)`,
}

func init() {
	RegisterScripts(`help`, Script{
		Action:  helpSetup,
		Ruleset: helpRules,
		Usage: map[string]string{
			"!help":          "查看所有脚本简介",
			"!help <script>": "查看script脚本的帮助信息",
		},
		Description: `显示已经注册的所有脚本简介或帮助信息`,
	})
}
