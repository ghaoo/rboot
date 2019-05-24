package rboot

import (
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

const (
	DefaultRobotName      = `Rboot`
	DefaultRobotProvider  = `cli`
	DefaultRobotMemorizer = `memory`

	DefaultHttpServerPort = `192.168.0.150:9900`
)

type Robot struct {
	name string

	signalChan chan os.Signal
	sync.Mutex
}

func New() *Robot {

	log.SetLevel(log.DebugLevel)

	bot := &Robot{}

	return bot
}

