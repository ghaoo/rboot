package rboot

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Rboot struct {
	input chan Message
	evtStream *evtStream
}

func NewBot() (*Rboot, error) {

	logrus.Info(`皮皮虾我们走...`)

	bot := &Rboot{
		evtStream: newEvtStream(),
	}

	bot.evtStream.init()

	bot.keepAlive()

	return bot, nil
}

// Go 皮皮虾我们走
func (bot *Rboot) Go() {
	es := bot.evtStream

	for e := range es.stream {
		switch e.Path {
		case "/sig/stoploop":
			logrus.Info(`皮皮虾快停下...`)
			return
		}
		go func(a Event) {
			es.RLock()
			defer es.RUnlock()
			if pattern := es.match(a.Path); pattern != "" {
				es.Handlers[pattern](a)
			}
		}(e)
		if es.hook != nil {
			es.hook(e)
		}
	}
}

func (bot *Rboot) process() {
	for in := range bot.input {

		go func(msg Message) {

			event := Event{
				Type: `NewMessage`,
				From: `Server`,
				Path: `/msg`,
				To:   `End`,
				Time: time.Now().Unix(),
				Data: msg,
			}
			bot.evtStream.serverEvt <- event
		}(in)
	}
}

func (bot *Rboot) keepAlive() {
	go func() {
		bot.process()
		bot.keepAlive()
		return
	}()
}

func (bot *Rboot) initialize() {
	logrus.SetLevel(logrus.DebugLevel)

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	// 接收消息
	// 指定消息提供者，如果配置文件没有指定，则默认使用 cli
	adpF, err := DetectAdapter(`cli`)

	if err != nil {
		panic(`Detect adapter error: ` + err.Error())
	}

	adp := adpF(bot)

	bot.input = adp.Incoming()

}