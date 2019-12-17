package rboot

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Rboot struct {
	Adapter Adapter
	Match     string
	Rule      Rule
	input     chan Message
	evtStream *evtStream


}

func NewBot() (*Rboot, error) {

	logrus.Info(`皮皮虾我们走...`)

	bot := &Rboot{
		evtStream: newEvtStream(),
	}

	bot.Rule = &Regex{}

	bot.evtStream.init()

	bot.initialize()

	bot.keepAlive()

	return bot, nil
}

// Go 皮皮虾我们走
func (bot *Rboot) Go() {
	es := bot.evtStream

	for k := range es.Handlers {
		logrus.Debugf(k)
	}

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

		go func(bot *Rboot, msg Message) {
			/*defer func() {
				if r := recover(); r != nil {
					logrus.Errorf("panic recovered when parsing message: %#v. Panic: %v", msg, r)
				}
			}()*/

			event := Event{
				Type: `NewMessage`,
				From: `Server`,
				Path: `/msg`,
				To:   `End`,
				Time: time.Now().Unix(),
				Data: msg,
			}
			bot.evtStream.serverEvt <- event

			if script, match, ok := bot.MatchRuleset(msg.Content); ok {

				bot.Match = match

				action, err := DirectiveScript(script)

				if err != nil {
					logrus.Error(err)
				}

				responses := action(bot)

				for _, resp := range responses {
					resp.From = msg.To
					resp.To = msg.From

					bot.Adapter.Send(resp)
				}
			}

		}(bot, in)
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
