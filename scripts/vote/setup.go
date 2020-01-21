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
		return vote.New(args[1], args[2])
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
		Usage: map[string]string{
			"!投票<name> [opt1 opt2...]": "新建投票",
			"@@<option>":               "为某一选项投票",
			"!投票结果":                    "查看投票结果",
			"!投票结束":                    "结束投票并返回结果",
		},
		Description: `投票小插件。`,
	})
}
