package scripts

import (
	"strings"

	"github.com/ghaoo/rboot"
	"log"
)

func setup(bot rboot.Robot, msg rboot.Message) []rboot.Message {

	if bot.MatchMessage(`ping|PING|Ping`, msg) {
		reg := bot.Regexp(`ping|PING|Ping`)

		fs := reg.FindAllStringSubmatch(msg.Content(), -1)

		log.Printf(`%v`, fs)

		return []rboot.Message{
			{
				Body: strings.NewReader(`PONG`),
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
