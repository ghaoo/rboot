package scripts

import (
	"github.com/ghaoo/rboot"
	"regexp"
	"time"
	"math/rand"
	"strconv"
)

func setup(bot rboot.Robot, msg rboot.Message) []rboot.Message {
	var reg *regexp.Regexp
	reg = regexp.MustCompile(`ping|PING`)

	if reg.MatchString(msg.Content) {
		bot.Send(rboot.Message{Content:randReply()})
	}

	return nil
}

func call(bot rboot.Robot) error {
	bot.Ticker(60 * time.Second)
	bot.Handle(`/ticker/1m`, func(evt rboot.Event) {
		data := evt.Data.(rboot.TickerData)

		bot.Send(rboot.Message{Content:`This is a minute-long task: PONG..., the ` + strconv.Itoa(int(data.Count)) + `th time`})
	})
	return nil
}

func randReply() string {
	now := time.Now()
	rand.Seed(int64(now.Unix()))
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
	return replies[rand.Intn(len(replies)-1)]
}

func init() {
	rboot.RegisterScript(`ping`, &rboot.Script{
		Action: setup,
		//Call: call,
	})
}
