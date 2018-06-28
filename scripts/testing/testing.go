package testing

import (
	"fmt"
	"rboot"
	"time"
)

func parse(bot rboot.Robot, msg rboot.Message) []rboot.Message {

	if msg.Content == `11` {
		return []rboot.Message{
			{Content: `22222222222222222`},
		}
	}
	if msg.Content == `2` {
		return []rboot.Message{
			{Content: `11111111111111111`},
		}
	}

	return nil
}

func hook(bot rboot.Robot) {

	bot.Ticker(2 * time.Second)
	bot.Handle(`/ticker/2s`, func(evt rboot.Event) {
		data := evt.Data.(rboot.TimerData)

		str := fmt.Sprintf(
			`
time: %v
count: %d
`, time.Now().Local(), data.Count)

		bot.Send(rboot.Message{Content: str})
	})

	bot.Timing(`15:14`)
	bot.Handle(`/timing/15:13`, func(evt rboot.Event) {

		bot.Send(rboot.Message{Content: time.Now().Local().String()})
	})
}

func init() {
	rboot.RegisterScript(`testing`, &rboot.Script{
		Action: parse,
		Hook:   hook,
	})
}
