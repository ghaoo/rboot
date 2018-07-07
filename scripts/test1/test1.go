package test1

import (
	"fmt"
	"github.com/ghaoo/rboot"
	"regexp"
)

func setup(bot rboot.Robot, msg rboot.Message) []rboot.Message {
	var reg *regexp.Regexp
	reg = regexp.MustCompile(`test1`)

	if reg.MatchString(msg.Content) {
		fmt.Println(`test1`)
		bot.Send(rboot.Message{Content: `test1`})
	}

	return nil
}

func init() {
	rboot.RegisterScript(`test1`, &rboot.Script{
		Action: setup,
		//Call: call,
	})
}
