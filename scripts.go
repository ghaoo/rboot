package rboot

import (
	"context"
	"fmt"
	"strings"
)

var (
	scripts  = make(map[string]Script)
	rulesets = make(map[string]map[string]string)
)

// Script 脚本结构体
type Script struct {
	Action      SetupFunc         // 执行解析或一些必要加载
	Ruleset     map[string]string // 脚本规则集合
	Usage       string            // 帮助信息
	Description string            // 简介
}

// SetupFunc 脚本执行或解析
type SetupFunc func(context.Context, *Robot) []Message

// RegisterScripts 注册脚本
func RegisterScripts(name string, script Script) {

	if name == "" {
		panic("RegisterScripts: the script must have a name")
	}
	if _, ok := scripts[name]; ok {
		panic("RegisterScripts: script named " + name + " already registered.")
	}

	scripts[name] = script

	if len(script.Ruleset) > 0 {

		rulesets[name] = script.Ruleset
	}
}

// DirectiveScript 根据脚本名称获取脚本执行函数
func DirectiveScript(name string) (SetupFunc, error) {

	if script, ok := scripts[name]; ok {
		return script.Action, nil
	}

	return nil, fmt.Errorf("DirectiveScript: no action found in script '%s' (missing a script?)", name)
}

// helpSetup 帮助脚本
func helpSetup(ctx context.Context, bot *Robot) []Message {

	switch bot.MatchRule {
	case `help`:
		if len(bot.MatchSub) < 2 {
			return []Message{
				{
					Content: "请在 !help 后面带上想要查看的脚本名称，比如查看 <ping> 脚本帮助信息，输入 <!help ping>",
				},
			}
		} else {
			return help(bot.MatchSub[1])
		}
	case `script`:
		// 获取所有脚本信息
		content := ""

		for scr, spt := range scripts {
			content += fmt.Sprintf("%s: %s", scr, spt.Description)
			content += "\n"
		}

		// 去除末尾空白字符
		content = strings.TrimSpace(content)

		return []Message{{Content: content}}
	}

	return nil
}

// help 帮助信息
func help(scr string) []Message {
	if script, ok := scripts[scr]; ok {

		return []Message{{Content: script.Usage}}
	} else {
		return []Message{{Content: "help命令用法：!help <script> \n!scripts 可查看所有加载的脚本信息"}}
	}

	return nil
}

// 帮助脚本规则集
var helpRules = map[string]string{
	`help`:   `^!help(?: *)(.*)`,
	`script`: `^!(?:脚本|scripts)`,
}

func init() {
	RegisterScripts(`help`, Script{
		Action:      helpSetup,
		Ruleset:     helpRules,
		Usage:       "!script 或 !脚本: 查看所有脚本 \n!help <script>: 查看脚本帮助信息",
		Description: `查看脚本信息`,
	})
}
