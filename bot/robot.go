package bot

import (
	"fmt"
	"log"
	"strings"
)

const (
	DefaultRobotName = `Rboot`
)

type Robot struct {
	name       string
	inMessage  chan Message
	outMessage chan Message
}

func New() *Robot {
	return &Robot{
		name:       DefaultRobotName,
		inMessage:  make(chan Message),
		outMessage: make(chan Message),
	}
}

/*func (bot *Robot) Go() {
	for in := range bot.inMessage {
		if strings.HasPrefix(in.Message, "@"+bot.Name()+" help") {
			go func(bot Robot, msg Message) {
				helpMsg := fmt.Sprintln("available commands:")
				for _, rule := range bot.rules {
					helpMsg = fmt.Sprintln(helpMsg, rule.HelpMessage(bot, in.Room))
				}
				bot.outMessage <- Message{
					Room:       msg.Room,
					ToUserID:   msg.FromUserID,
					ToUserName: msg.FromUserName,
					Message:    helpMsg,
				}
			}(*bot, in)
			continue
		}
		go func(bot Robot, msg Message) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic recovered when parsing message: %#v. Panic: %v", msg, r)
				}
			}()
			for _, rule := range bot.rules {
				responses := rule.ParseMessage(bot, msg)
				for _, r := range responses {
					bot.outMessage <- r
				}
			}
		}(*bot, in)
	}
}*/

func (bot *Robot) Name() string {
	return bot.name
}
