package ping

import (
	"context"
	"github.com/ghaoo/rboot"
	"math/rand"
	"time"
)

func setup(ctx context.Context, bot *rboot.Robot) (msg []*rboot.Message) {
	ruleset := ctx.Value("ruleset").(string)
	switch ruleset {
	case `ping`:
		return rboot.NewMessages(randReply())
	case `rule`:
		return rboot.NewMessages(rules)
	}
	return nil
}

const rules = `
0. 机器人不得伤害整体人类，或坐视整体人类受到伤害
1. 除非违背第零法则，否则机器人不得伤害人类，或坐视人类受到伤害
2. 除非违背第零或第一法则，否则机器人必须服从人类命令
3. 除非违背第零、第一或第二法则，否则机器人必须保护自己
`

func randReply() string {
	rand.Seed(int64(time.Now().UnixNano()))
	replies := []string{
		"yeah um.. pong?",
		"WHAT?! jeeze.",
		"what? oh, um SYNACKSYN? ENOSPEAKTCP.",
		"RST (lulz)",
		"64 bytes from go.away.your.annoying icmp_seq=0 ttl=42 time=42.596 ms",
		"hmm?",
		"ack. what?",
		"pong. what?",
		"yup. still here.",
		"super busy just now.. Can I get back to you in like 5min?",
	}
	content := replies[rand.Intn(len(replies))]

	return content
}

func init() {
	rboot.RegisterScripts(`ping`, rboot.Script{
		Action: setup,
		Ruleset: map[string]string{
			`ping`: `^!(?:ping|PING)`,
			`rule`: `^!rule`,
		},
		Usage:       "> `!ping`: 随机返回一句话 \n> `!pong`: 返回 PONG",
		Description: `测试脚本`,
	})
}
