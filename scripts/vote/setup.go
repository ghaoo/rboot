package vote

import (
	"github.com/ghaoo/rboot"
)

var vote = new(Vote)

func setup(bot *rboot.Robot, in *rboot.Message) []*rboot.Message {
	rule := in.Header.Get("rule")
	args := in.Header["args"]

	switch rule {
	case `new_vote`:
		return vote.New(bot, args[1], args[2])
	case `voting`:
		return vote.Voting(in.Sender, args[1])
	case `stop_vote`:
		return vote.Stop()
	case `result`:
		return vote.Result()
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
		Usage: "> `!投票<投票名称> [选项1 选项2...]`: 新建投票\n\n" +
			"> `@@<选项>`: 投票某一选项\n\n " +
			"> `!投票结果`: 查看投票结果\n\n " +
			"> `!投票结束`: 结束投票",
		Description: `投票小插件。`,
	})
}
