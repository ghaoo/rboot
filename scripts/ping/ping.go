package ping

import (
	"github.com/ghaoo/rboot"
	"math/rand"
	"time"
)

func setup(bot *rboot.Robot, in *rboot.Message) (msg []*rboot.Message) {
	return rboot.NewMessages(randReply())
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
			`ping`: `^!(?:ping|PING)`,
			`rule`: `^!rule`,
		},
		Usage:       "> `!ping`: 随机返回一句话 \n\n> `!pong`: 返回 PONG",
		Description: `测试脚本`,
	})
}
