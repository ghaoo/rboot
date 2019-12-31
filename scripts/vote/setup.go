package vote

import (
	"context"
	"github.com/ghaoo/rboot"
)

var vote = new(Vote)

func setup(ctx context.Context, bot *rboot.Robot) []rboot.Message {
	in := ctx.Value("input").(rboot.Message)

	/*if in.Mate["GroupMsg"] != nil && !in.Mate["GroupMsg"].(bool) {
		return []rboot.Message{{Content: "投票请在群组中创建"}}
	}*/

	switch bot.Ruleset {
	case `vote`:
		return vote.Voting(in.Sender.Name, bot.Args[1])
	case `new_vote`:
		return vote.New(bot, in.From, bot.Args[1], in.Sender.Name, bot.Args[2])
	case `stop_vote`:
		return vote.Stop(in.Sender.Name)
	case `result`:
		return vote.Result()
	}

	return nil
}

func init() {
	rboot.RegisterScripts(`vote`, rboot.Script{
		Action: setup,
		Ruleset: map[string]string{
			`vote`:      `^@@(.+)`,
			`result`:    `^!投票结果`,
			`new_vote`:  `^!投票(.+)[ ]?\[(.+)\]`,
			`stop_vote`: `^!结束投票`,
		},
		Usage:       "",
		Description: `投票小插件。`,
	})
}
