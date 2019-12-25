package help

import (
	"context"
	"fmt"
	"github.com/ghaoo/rboot"
	"strings"
)

var scripts = rboot.ListScripts()

// helpSetup 帮助脚本
func helpSetup(ctx context.Context, bot *rboot.Robot) []rboot.Message {

	switch bot.MatchRule {
	case `help`:
		if len(bot.MatchSub) < 2 {
			return []rboot.Message{
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

		return []rboot.Message{{Content: content}}
	}

	return nil
}

// help 帮助信息
func help(scr string) []rboot.Message {
	if script, ok := scripts[scr]; ok {

		return []rboot.Message{{Content: script.Usage}}
	} else {
		return []rboot.Message{{Content: "help命令用法：!help <script> \n!scripts 可查看所有加载的脚本信息"}}
	}

	return nil
}

// 帮助脚本规则集
var helpRules = map[string]string{
	`help`:   `^!help(?: *)(.*)`,
	`script`: `^!(?:脚本|scripts)`,
}

func init() {
	rboot.RegisterScripts(`help`, rboot.Script{
		Action:      helpSetup,
		Ruleset:     helpRules,
		Usage:       "!script 或 !脚本: 查看所有脚本 \n!help <script>: 查看脚本帮助信息",
		Description: `查看脚本信息`,
	})
}
