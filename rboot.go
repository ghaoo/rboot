package rboot

import (
	"sync"
	"log"
)

type Rboot struct {
	name string
	providerIn  chan Message
	providerOut chan Message

	plugs map[string]Plugin
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

			go func(bot Rboot, msg Message) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("panic recovered when parsing message: %#v. Panic: %v", msg, r)
					}
				}()
				/*for _, plug := range bot.plugs {
					responses := plug.Action
					for _, r := range responses {
						s.providerOut <- r
					}
				}*/
			}(*bot, in)
		}
	})
}

func (bot *Rboot) Name() string {
	return bot.name
}

// 适配连接器
type Provider interface {
	Run() error
	Incoming() chan Message
	Outgoing() chan Message
}
