package scripts

import (
	"github.com/ghaoo/rboot"
	"regexp"
	"log"
	"time"
)

func setup(bot rboot.Robot, msg rboot.Message) []rboot.Message {
	var reg *regexp.Regexp
	reg = regexp.MustCompile(`1|2|3|4`)

	var match []string

	if reg.MatchString(msg.Content) {
		match = reg.FindAllStringSubmatch(msg.Content, -1)[0]
	}

	log.Printf(`%v`, match)

	return nil
}

func call(bot rboot.Robot) error {
	bot.Ticker(2 * time.Second)
	bot.Handle(`/ticker/2s`, func(evt rboot.Event) {
		//data := evt.Data.(rboot.TimerData)

		bot.Send(rboot.NewStringMessage(`111111111`))
	})
	return nil
}


func init() {
	rboot.RegisterScript(`ping`, &rboot.Script{
		Action: setup,
		Call: call,
	})
}
