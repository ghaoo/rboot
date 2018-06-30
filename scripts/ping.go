package scripts

import (
	"strings"

	"github.com/ghaoo/rboot"
	"log"
)

func setup(bot rboot.Robot, msg rboot.Message) []rboot.Message {

	println(msg.Content())

	if bot.MatchMessage(`ping|PING|Ping`, msg) {
		println(msg.Content())
		reg := bot.Regexp(`ping|PING|Ping`)

		fs := reg.FindAllStringSubmatch(msg.Content(), -1)

		log.Printf(
			`msg: %v
fs1: %v
fs2: %v
`,msg.Body, fs, reg.FindAllString(msg.Content(), -1))

		return []rboot.Message{
			{
				Body: strings.NewReader(`PONG ......`),
			},
		}
	}

	return nil
}

func init() {
	rboot.RegisterScript(`ping`, &rboot.Script{
		Action: setup,
	})
}
