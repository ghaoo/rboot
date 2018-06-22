package rboot

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	DefaultRbootConf      = `config.yml`
	DefaultRobotName      = `Rboot`
	DefaultRobotProvider = `cli`
)

type Rboot struct {
	name      string
	provider Provider
	conf      Config

	signalChan chan os.Signal
}

func NewRboot(config ...string) *Rboot {

	var conf = DefaultRbootConf

	if len(config) > 0 {
		conf = config[0]
	}

	bot := &Rboot{
		conf:       NewConf(conf),
		signalChan: make(chan os.Signal, 1),
	}

	return bot
}

func (bot *Rboot) SetName(name string) {
	bot.name = name
}

func (bot *Rboot) SetProvider(provider Provider) {
	bot.provider = provider
}

func (bot *Rboot) Conf() Config {
	return bot.conf
}

func (bot *Rboot) Go() {
	bot.initialize()

	go bot.provider.Run()

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

func (bot *Rboot) Stop() error {

	log.Printf("stopping %s connecter", bot.provider.Name())
	if err := bot.provider.Close(); err != nil {
		return err
	}

	log.Printf("stopping %s", DefaultRobotName)
	return nil
}

func (bot *Rboot) Name() string {
	return bot.name
}

func (bot *Rboot) initialize() {
	if bot.conf.Name == `` {
		bot.name = DefaultRobotName
	} else {
		bot.name = bot.conf.Name
	}

	res := NewResponse(bot)
	botConName := DefaultRobotProvider

	if bot.conf.Connecter != `` {
		botConName = bot.conf.Connecter
	}

	con, err := getProvider(res, botConName)

	if err != nil {
		panic(`initialize error: ` + err.Error())
	}

	bot.provider = con
}
