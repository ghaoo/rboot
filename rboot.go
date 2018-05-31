package rboot

import (
	"regexp"
	"sync"
	"log"
)

type Rboot struct {
	name        string
	providerIn  chan Message
	providerOut chan Message
}

func NewRboot(name string) *Rboot {
	return &Rboot{
		name:        name,
		providerIn:  make(chan Message),
		providerOut: make(chan Message),
	}
}

var once sync.Once

func (bot *Rboot) Go() {
	once.Do(func() {

		for in := range bot.providerIn {

			go func(bot *Rboot, msg Message) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Rboot: panic recovered when parsing message: %#v. Panic: %v", msg, r)
					}
				}()
				mbyte, err := msg.Read()
				if err != nil {
					panic(err)
				}

				for name, rules := range ListPlugins() {
					if bot.match(rules, string(mbyte)) {
						handle, err := DirectiveAction(name)
						if err != nil {
							panic(err)
						}

						c := &Controller{bot}

						outs, err := handle(c)

						for _, out := range outs {
							bot.providerOut <- out
						}
					}
				}
			}(bot, in)
		}
	})
}

func (bot *Rboot) regexp(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func (bot *Rboot) match(rules []string, msg string) bool {

	for _, rule := range rules {
		reg := bot.regexp(rule)

		if !reg.MatchString(msg) {
			return false
		}
		return true
	}
	return false
}

func (bot *Rboot) Name() string {
	return bot.name
}

func (bot *Rboot) Incoming() chan Message {
	return bot.providerIn
}

func (bot *Rboot) Outgoing() chan Message {
	return bot.providerOut
}
