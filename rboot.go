package rboot

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	DefaultRobotName = `Rboot`
)

type Rboot struct {
	name      string
	connecter Connecter

	signalChan chan os.Signal
}

func NewRboot() *Rboot {
	bot := &Rboot{
		name:       DefaultRobotName,
		signalChan: make(chan os.Signal, 1),
	}

	return bot
}

func (bot *Rboot) SetName(name string) {
	bot.name = name
}

func (bot *Rboot) SetConnecter(connecter Connecter) {
	bot.connecter = connecter
}

func (bot *Rboot) Go() {

	go bot.connecter.Run()

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

	log.Printf("stopping %s connecter", bot.connecter.Name())
	if err := bot.connecter.Close(); err != nil {
		return err
	}

	log.Printf("stopping %s", DefaultRobotName)
	return nil
}

func (bot *Rboot) Name() string {
	return bot.name
}
