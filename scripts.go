package rboot

import (
	"fmt"
	"strings"
)

var (
	scripts = make(map[string]Script)
	ruleset = make(map[string]map[string]string)
)

// Script 脚本结构体
type Script struct {
	Action      ScriptFunc        // 执行解析或一些必要加载
	Ruleset     map[string]string // 脚本规则集合
	Usage       string            // 帮助信息
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
			if script, ok := scripts[args[1]]; ok {
				msg = append(msg, NewMessage(script.Usage))
			} else {
				msg = append(msg, NewMessage("> help命令用法：!help <script> \n\n> !scripts 可查看所有加载的脚本信息"))
			}
		}
	case `ruleset`:
		if len(args) < 2 || args[1] == "" {
			content := ""
			for scr, spt := range scripts {
				content += fmt.Sprintf("**%s**:\n", scr)
				for ruleset := range spt.Ruleset {
					content += fmt.Sprintf("- %s\n", ruleset)
				}

				content += "\n"
			}

			content = strings.TrimSpace(content)

			msg = append(msg, NewMessage(content))

		} else {

			scr := args[1]
			spt := scripts[scr]
			content := fmt.Sprintf("**%s**:\n", scr)

			for ruleset := range spt.Ruleset {
				content += fmt.Sprintf("- %s\n", ruleset)
			}

			msg = append(msg, NewMessage(content))
		}
	}

	return msg
}

func script() string {
	// 获取所有脚本信息
	content := ""

	for scr, spt := range scripts {
		content += fmt.Sprintf(" **%s**: %s\n **Usage**:\n%s", scr, spt.Description, spt.Usage)
		content += "\n\n"
	}

	// 去除末尾空白字符
	content = strings.TrimSpace(content)

	return content
}

// 帮助脚本规则集
var helpRules = map[string]string{
	`help`:    `^!help(?: *)(\S*)`,
	`ruleset`: `^!ruleset(?: *)(\S*)`,
}

func init() {
	RegisterScripts(`help`, Script{
		Action:      helpSetup,
		Ruleset:     helpRules,
		Usage:       "> `!help <script>`: 查看脚本帮助信息 \n\n> `!ruleset <script>`: 查看已经注册的脚本规则集",
		Description: `查看脚本帮助信息`,
	})
}
