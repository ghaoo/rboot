package rboot

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
)

const (
	DefaultRbootConf     = `config.yml`
	DefaultRobotName     = `Rboot`
	DefaultRobotProvider = `cli`
)

type Robot struct {
	name     string
	es       *eventStream
	provider Provider
	conf     Config

	Matcher string

	sync.Mutex
	signalChan chan os.Signal
}

func NewRboot(config ...string) *Robot {

	var conf = DefaultRbootConf

	if len(config) > 0 {
		conf = config[0]
	}

	bot := &Robot{
		es:         newStream(),
		conf:       NewConf(conf),
		signalChan: make(chan os.Signal, 1),
	}

	return bot
}

func (bot *Robot) SetName(name string) {
	bot.name = name
}

func (bot *Robot) SetProvider(provider Provider) {
	bot.provider = provider
}

func (bot *Robot) Conf() Config {
	return bot.conf
}

// 皮皮虾，我们走~~~~~~~~~
func (bot *Robot) Go() {
	bot.initialize()

	go bot.provider.Run()

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

	log.Printf("stopping %s connecter", bot.provider.Name())
	if err := bot.provider.Close(); err != nil {
		return err
	}

	log.Printf("stopping %s", DefaultRobotName)
	return nil
}

func (bot *Robot) Name() string {
	return bot.name
}

func (bot *Robot) Receive(msg *Message) error {
	bot.Lock()
	defer bot.Unlock()

	b, err := msg.Read()
	if err != nil {
		return fmt.Errorf(`Receive: message read error %v `, err)
	}

	if msg.Header.From() == `` {
		msg.Header.Set(`From`, `System`)
	}

	if msg.Header.To() == `` {
		msg.Header.Set(`To`, `Nil`)
	}

	text := string(b)

	scrName, ok := bot.matchRuleset(text)

	if ok {

		action, err := DirectiveAction(scrName)

		if err != nil {
			return err
		}

		return action(bot)
	}

	return fmt.Errorf(`Receive: no matching scripts... `)
}

func (bot *Robot) ReceiveWithReader(in io.Reader) error {
	msg, err := ReadMessage(in)

	if err != nil {
		return err
	}

	return bot.Receive(msg)
}

func (bot *Robot) matchRuleset(msg string) (string, bool) {
	for scr, rules := range rulesets {
		for matcher, rule := range rules {
			if bot.match(rule, msg) {
				bot.Matcher = matcher
				return scr, true
			}
		}
	}

	log.Printf(`no match script`)
	return ``, false
}

func (bot *Robot) match(pattern, msg string) bool {

	reg := regexp.MustCompile(pattern)

	if reg.MatchString(msg) {
		return true
	}

	return false
}

func (bot *Robot) Send(strs ...string) error {
	return bot.provider.Send(strs...)
}

func (bot *Robot) Reply(strs ...string) error {
	return bot.provider.Reply(strs...)
}

func (bot *Robot) initialize() {

	if bot.conf.Name == `` {
		bot.name = DefaultRobotName
	} else {
		bot.name = bot.conf.Name
	}

	botConName := DefaultRobotProvider

	if bot.conf.Connecter != `` {
		botConName = bot.conf.Connecter
	}

	con, err := getProvider(bot, botConName)

	if err != nil {
		panic(`initialize error: ` + err.Error())
	}

	bot.provider = con

	bot.es.init()

	bot.es.merge("custom", usrEvent)
}
