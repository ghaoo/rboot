package rboot

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"regexp"
)

const (
	DefaultRbootConf     = `config.yml`
	DefaultRobotName     = `Rboot`
	DefaultRobotProvider = `cli`
)

type Robot struct {
	name string
	es   *eventStream
	conf Config

	providerIn  chan Message
	providerOut chan Message

	signalChan chan os.Signal
	sync.Mutex
}

func New(confpath ...string) *Robot {

	var conf = DefaultRbootConf

	if len(confpath) > 0 {
		conf = confpath[0]
	}

	bot := &Robot{
		es:          newStream(),
		conf:        NewConf(conf),
		providerIn:  make(chan Message),
		providerOut: make(chan Message),
		signalChan:  make(chan os.Signal, 1),
	}

	return bot
}

func (bot *Robot) SetName(name string) {
	bot.name = name
}

func (bot *Robot) Conf() Config {
	return bot.conf
}

func (bot *Robot) Send(msg Message) {
	bot.providerOut <- msg
}

var processOnce sync.Once

func (bot *Robot) process() {
	processOnce.Do(func() {

		for _, script := range availableScripts {
			script.Hook(*bot)
		}

		for in := range bot.providerIn {
			go func(bot Robot, msg Message) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("panic recovered when parsing message: %#v. Panic: %v", msg, r)
					}
				}()

				for _, script := range availableScripts {
					responses := script.Action(bot, in)

					for _, r := range responses {
						bot.providerOut <- r
					}
				}

			}(*bot, in)
		}
	})
}

// 皮皮虾，我们走~~~~~~~~~
func (bot *Robot) Go() {
	bot.initialize()

	go bot.process()

	go bot.es.loop()

	signal.Notify(bot.signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	stop := false
	for !stop {
		select {
		case sig := <-bot.signalChan:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				stop = true
			}
		}
	}

	signal.Stop(bot.signalChan)

	bot.Stop()
}

func (bot *Robot) Stop() error {

	log.Printf("stopping %s", DefaultRobotName)
	return nil
}

func (bot *Robot) Name() string {
	return bot.name
}

func (bot *Robot) initialize() {

	if bot.conf.Name == `` {
		bot.name = DefaultRobotName
	} else {
		bot.name = bot.conf.Name
	}

	bot.es.init()

	bot.es.merge("custom", usrEvent)

	// 指定消息提供者，如果配置文件没有指定，则默认使用 cli
	provName := DefaultRobotProvider

	if bot.conf.Provider != `` {
		provName = bot.conf.Provider
	}

	prov, err := Detect(provName)

	if err != nil {
		panic(`Detect error: ` + err.Error())
	}

	bot.providerIn = prov.Incoming()
	bot.providerOut = prov.Outgoing()
}

func (bot *Robot) Regexp(pattern string) *regexp.Regexp {
	return regex(pattern)
}

func (bot *Robot) MatchMessage(pattern string, msg Message) bool {
	return match(pattern, msg)
}

func regex(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func match(pattern string, msg Message) bool {

	if !regex(pattern).MatchString(msg.Content()) {
		return false
	}

	return true

}