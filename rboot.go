package rboot

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	DefaultRbootConf = `config.yml`
	DefaultRobotName = `Rboot`
	DefaultRobotProvider = `cli`
)

type Robot struct {
	name    string
	es      *eventStream
	conf    Config

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

// 皮皮虾，我们走~~~~~~~~~
func (bot *Robot) Go() {
	bot.initialize()

	go bot.Process()

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

func (bot *Robot) Send(msg Message) {
	bot.providerOut <- msg
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
