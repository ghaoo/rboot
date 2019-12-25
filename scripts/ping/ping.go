package ping

import (
	"context"
	"github.com/ghaoo/rboot"
	"math/rand"
	"time"
)

func setup(ctx context.Context, bot *rboot.Robot) []rboot.Message {
	var msg []rboot.Message

	switch bot.MatchRule {
	case `ping`:

		msg = []rboot.Message{
			{
				Content: randReply(),
			},
		}
	case `pong`:

		msg = []rboot.Message{
			{
				Content: `Pong!`,
			},
		}

	case `history`:
		msg = []rboot.Message {
			{
				Content: "" + bot.PrevHistory("").Incoming().Content,
			},
		}
	}

	return msg
}

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
			`ping`: `^!(ping|PING)`,
			`pong`: `^!(pong|PONG)`,
			`history`: `^!(history|历史)`,
		},
		Usage:       "!ping: 随机返回一句话 \n!pong: 返回 PONG",
		Description: `测试程序和脚本是否运行正常。`,
	})
}
