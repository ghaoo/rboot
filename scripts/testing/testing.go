package testing

import (
	"fmt"
	"time"

	"github.com/ghaoo/rboot"
	"strings"
)

func setup(bot rboot.Robot, msg rboot.Message) []rboot.Message {

	/*if msg.Content == `11` {
		return []rboot.Message{
			{Content: `22222222222222222`},
		}
	}
	if msg.Content == `2` {
		return []rboot.Message{
			{Content: `11111111111111111`},
		}
	}*/

	return nil
}

func hook(bot rboot.Robot) {

	bot.Ticker(63 * time.Second)
	bot.Handle(`/ticker/1m3s`, func(evt rboot.Event) {
		data := evt.Data.(rboot.TickerData)

		str := fmt.Sprintf(
			`
time: %v
count: %d
`, time.Now().Local(), data.Count)

		bot.Send(rboot.Message{Body: strings.NewReader(str)})
	})

	bot.Timing(`15:14`)
	bot.Handle(`/timing/15:13`, func(evt rboot.Event) {

		bot.Send(rboot.Message{Body: strings.NewReader(time.Now().Local().String())})
	})
}

func init() {
	rboot.RegisterScript(`testing`, &rboot.Script{
		Action: setup,
		Hook:   hook,
	})
}
