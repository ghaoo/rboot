package help

import (
	"fmt"
	"strings"
	"github.com/ghaoo/rboot"
	"context"
)

var scripts = rboot.GetScripts()

func setup(ctx context.Context, bot *rboot.Robot) []rboot.Message {

	switch bot.Match {
	case `help`:
		if len(bot.MatchString) < 2 {
			return []rboot.Message{
				{
					Content: "请在 !help 后面带上想要查看的脚本名称，比如查看 <ping> 脚本帮助信息，输入 <!help ping>",
				},
			}
		} else {
			return help(bot.MatchString[1])
		}
	case `script`:
		return getScript()
	}

	return nil
}

func getScript() []rboot.Message {
	scrs := ""

	for scr, spt := range scripts {
		scrs += fmt.Sprintf("%s: %s", scr, spt.Description)
		scrs += "\n"
	}

	scrs = strings.TrimSpace(scrs)

	return []rboot.Message{{Content: scrs}}
}

func help(scr string) []rboot.Message {
	if script, ok := scripts[scr]; ok {

		return []rboot.Message{{Content: script.Usage}}
	} else {
		return []rboot.Message{{Content: "help命令用法：!help <script>"}}
	}

	return nil
}

var helpRules = map[string]string{
	`help`:   `^!help(?: +)(.*)`,
	`script`: `^!(?:脚本|scripts)`,
}

func init() {
	rboot.RegisterScripts(`help`, rboot.Script{
		Action:      setup,
		Ruleset:     helpRules,
		Usage:       "!script 或 !脚本: 查看所有脚本 \n!help <script>: 查看脚本帮助信息",
		Description: `查看脚本信息`,
	})
}
