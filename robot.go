package rboot

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"sync"
	"context"
	"github.com/ghaoo/rboot/tools/env"
)

var AppName string

const Version = "3.0.0"

type Robot struct {
	Memory      Memorizer
	Match       string
	MatchString []string
	Rule        Rule
	Adapter     Adapter
	Contacts    []User
	conf        Config

	inputChan  chan Message
	outputChan chan Message

	sync.RWMutex
}

func New() *Robot {

	bot := &Robot{
		inputChan:  make(chan Message),
		outputChan: make(chan Message),
		conf:       newConfig(),
		Rule:       new(Regex),
	}

	return bot
}

var processOnce sync.Once

func process(ctx context.Context, bot *Robot) {
	processOnce.Do(func() {

		for in := range bot.inputChan {
			go func(bot Robot, msg Message) {
				defer func() {
					if r := recover(); r != nil {
						logrus.Errorf("panic recovered when parsing message: %#v. Panic: %v", msg, r)
					}
				}()

				ctx = context.WithValue(ctx, "input", msg)

				if script, match, sub, ok := bot.MatchRuleset(msg.Content); ok {

					bot.Match = match

					bot.MatchString = sub

					action, err := DirectiveScript(script)

					if err != nil {
						logrus.Error(err)
					}

					responses := action(ctx, &bot)

					if len(responses) > 0 {
						for _, resp := range responses {
							resp.From = msg.To
							resp.To = msg.From

							bot.outputChan <- resp
						}
					}

				}

			}(*bot, in)

		}
	})

}

// 皮皮虾，我们走~~~~~~~~~
func (bot *Robot) Go() {

	logrus.Debugf("Rboot Version %s", Version)
	// 机器人名称
	AppName = bot.conf.Name

	bot.initialize()

	logrus.Debug(`皮皮虾，我们走~~~~~~~`)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = context.WithValue(ctx, "appname", AppName)
	process(ctx, bot)

	bot.Stop()
}

// 皮皮虾，快停下~~~~~~~~~
func (bot *Robot) Stop() error {

	runtime.SetFinalizer(bot, nil)

	logrus.Debug(`皮皮虾，快停下~~~`)

	os.Exit(0)

	return nil
}

func (bot *Robot) SyncUsers(user []User) {

	if len(user) > 0 {
		bot.Contacts = user
	}
}

func (bot *Robot) SetMemo(memo Memorizer) *Robot {
	bot.Memory = memo

	return bot
}

func (bot *Robot) Send(msg Message) {
	bot.outputChan <- msg
}

func (bot *Robot) SendText(text string, to User) {
	msg := Message{
		To:      to,
		Content: text,
	}
	bot.outputChan <- msg
}

func (bot *Robot) MatchRuleset(msg string) (plug, match string, substr []string, matched bool) {

	for plug, rule := range rulesets {
		for m, r := range rule {
			if sub, ok := bot.Rule.Match(r, msg); ok {
				return plug, m, sub, true
			}
		}
	}

	return ``, ``, nil, false
}

func (bot *Robot) initialize() {

	if len(adapters) > 0 {

	}
	// 指定消息提供者，如果配置文件没有指定，则默认使用 cli
	adp, err := DetectAdapter(bot.conf.Adapter)

	if err != nil {
		panic(`Detect adapter error: ` + err.Error())
	}

	adapter := adp(bot)

	bot.inputChan = adapter.Incoming()
	bot.outputChan = adapter.Outgoing()

	// 指定储存器
	memo, err := DetectMemorizer(bot.conf.Memorizer)

	if err != nil {
		logrus.Errorf(`Detect memorizer error: %v`, err)
	}

	bot.Memory = memo()

}

func init() {

	// 加载配置
	err := env.Load()

	if err != nil {
		panic(`Load env config error: ` + err.Error())
	}
}
