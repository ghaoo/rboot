package help

import (
	"context"
	"fmt"
	"github.com/ghaoo/rboot"
	"strings"
)

var scripts = rboot.ListScripts()

// helpSetup 帮助脚本
func helpSetup(ctx context.Context, bot *rboot.Robot) (msg *rboot.Message) {

	switch bot.Ruleset {
	case `help`:
		if len(bot.Args) < 2 || bot.Args[1] == "" {
			msg = rboot.NewMessage("请在 !help 后面带上想要查看的脚本名称，比如查看 <ping> 脚本帮助信息，输入 <!help ping>")
		} else {
			if script, ok := scripts[bot.Args[1]]; ok {
				msg = rboot.NewMessage(script.Usage)
			} else {
				msg = rboot.NewMessage("> help命令用法：!help <script> \n> !scripts 可查看所有加载的脚本信息")
			}
		}
	case `ruleset`:
		if len(bot.Args) < 2 || bot.Args[1] == "" {
			content := ""
			for scr, spt := range scripts {
				content += fmt.Sprintf("**%s**:\n", scr)
				for ruleset := range spt.Ruleset {
					content += fmt.Sprintf("- %s\n", ruleset)
				}

				content += "\n"
			}

			content = strings.TrimSpace(content)

			msg = rboot.NewMessage(content)

		} else {

			scr := bot.Args[1]
			spt := scripts[scr]
			content := fmt.Sprintf("**%s**:\n", scr)

			for ruleset := range spt.Ruleset {
				content += fmt.Sprintf("- %s\n", ruleset)
			}

			msg = rboot.NewMessage(content)
		}
	case `script`:
		// 获取所有脚本信息
		content := ""

		for scr, spt := range scripts {
			content += fmt.Sprintf("- **%s**: %s\n- **Usage**:\n%s", scr, spt.Description, spt.Usage)
			content += "\n\n"
		}

		// 去除末尾空白字符
		content = strings.TrimSpace(content)

		msg = rboot.NewMessage(content)
	}

	return msg
}

// 帮助脚本规则集
var helpRules = map[string]string{
	`help`:    `^!help(?: *)(\S*)`,
	`ruleset`: `^!ruleset(?: *)(\S*)`,
	`script`:  `^!(?:脚本|scripts)`,
}

func init() {
	rboot.RegisterScripts(`help`, rboot.Script{
		Action:      helpSetup,
		Ruleset:     helpRules,
		Usage:       "> `!scripts` 或 `!脚本`: 查看所有脚本 \n> `!help <script>`: 查看脚本帮助信息 \n> `!ruleset <脚本名称>`",
		Description: `查看脚本帮助信息`,
	})
}
