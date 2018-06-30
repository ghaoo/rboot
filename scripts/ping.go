package scripts

import (
	"github.com/ghaoo/rboot"
	"time"
)

func setup(bot *rboot.Robot) error {
	switch bot.Pattern {
	case `ping`:
		bot.Send(rboot.NewStringMessage(`PONG...`))
	}

	return nil
}

func call(bot *rboot.Robot) error {
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
		Ruleset: map[string]string{
			`ping`:`ping|PING|Ping`,
		},
		Call: call,
	})
}
