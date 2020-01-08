package vote

import (
	"context"
	"github.com/ghaoo/rboot"
)

var vote = new(Vote)

func setup(ctx context.Context, bot *rboot.Robot) *rboot.Message {
	in := ctx.Value("input").(*rboot.Message)

	switch bot.Ruleset {
	case `voting`:
		return vote.Voting(bot.GetUserName(in.Sender), bot.Args[1])
	case `new_vote`:
		return vote.New(bot, in, bot.Args[1], bot.Args[2])
	case `stop_vote`:
		return vote.Stop(bot, in.Sender, bot.GetUserName(in.From))
	case `result`:
		return vote.Result(in.From)
	}

	return nil
}

func init() {
	rboot.RegisterScripts(`vote`, rboot.Script{
		Action: setup,
		Ruleset: map[string]string{
			`voting`:    `^@@(.+)`,
			`result`:    `^!投票结果`,
			`new_vote`:  `^!投票(.+)[ ]?\[(.+)\]`,
			`stop_vote`: `^!投票结束`,
		},
		Usage: "> `!投票<投票名称> [选项1 选项2]...`: 新建投票\n" +
			"> `@@<选项>`: 投票某一选项\n " +
			"> `!投票结果`: 查看投票结果\n " +
			"> `!投票结束`: 结束投票",
		Description: `投票小插件。`,
	})
}
