package rboot

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"fmt"
	"regexp"
)

const (
	DefaultRobotName      = `Rboot`
	DefaultRobotProvider  = `cli`
	DefaultRobotMemorizer = `memory`
)

type Robot struct {
	name string
	es   *eventStream
	memo Memorizer
	prov Provider
	Pattern string

	inMessage  chan Message

	signalChan chan os.Signal
	sync.Mutex

	err chan error
}

func New() *Robot {

	bot := &Robot{
		es:          newStream(),
		inMessage:  make(chan Message),
		signalChan:  make(chan os.Signal, 1),
	}

	return bot
}

// robot name
func (bot *Robot) Name() string {
	return bot.name
}

func (bot *Robot) SetName(name string) {
	bot.name = name
}

func (bot *Robot) Incoming() chan Message {
	return bot.inMessage
}

func (bot *Robot) Send(msg ...Message) {
	bot.prov.Send(msg...)
}

func (bot *Robot) Reply(msg ...Message) {
	bot.prov.Reply(msg...)
}

func (bot *Robot) Error(err error) {
	bot.err <- err
}

var processOnce sync.Once

func (bot *Robot) process() {

	processOnce.Do(func() {

		go func(bot *Robot) {
			defer func() {
				if r := recover(); r != nil {
					bot.err <- fmt.Errorf("panic recovered when executing script call: %v", r)
				}
			}()
			for sname, call := range execCall {
				err := call(bot)

				if err != nil {
					bot.err <- fmt.Errorf(`executing script(%s) call error: %v`, sname, err)
				}
			}

			// 处理错误信息
			for err := range bot.err {
				log.Print(err)
			}
		}(bot)

		for in := range bot.inMessage {
			go func(bot *Robot, msg Message) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("panic recovered when parsing message: %#v. Panic: %v", msg, r)
					}
				}()

				err := bot.handleMessage(in)

				if err != nil {
					bot.err <- err
				}

			}(bot, in)
		}
	})
}

// 处理消息
func (bot *Robot) handleMessage(msg Message) error {
	bot.Lock()
	defer bot.Unlock()

	if msg.Header.From() == `` {
		msg.Header.Set(`From`, `System`)
	}

	if msg.Header.To() == `` {
		msg.Header.Set(`To`, `Nil`)
	}

	scrName, ok := bot.matchRuleset(msg.Content)

	if ok {

		action, err := DirectiveAction(scrName)

		if err != nil {
			return err
		}

		return action(bot)
	}

	return fmt.Errorf(`Receive: no matching scripts... `)
}

// 匹配脚本规则集
func (bot *Robot) matchRuleset(msg string) (string, bool) {
	for scr, rules := range rulesets {
		for pattern, rule := range rules {
			if bot.match(rule, msg) {
				bot.Pattern = pattern
				return scr, true
			}
		}
	}

	return ``, false
}

// 匹配消息
func (bot *Robot) match(pattern, msg string) bool {

	reg := regexp.MustCompile(pattern)

	if reg.MatchString(msg) {
		return true
	}

	return false
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

// 皮皮虾，快停下~~~~~~~~~
func (bot *Robot) Stop() error {

	log.Printf(`stopping provider %s`, bot.prov.Name())
	bot.prov.Close()

	log.Printf("stopping %s", bot.name)
	return nil
}

// memorizer save data
func (bot *Robot) MemoSave(key string, value []byte) {
	bot.memo.Save(key, value)
}

// memorizer read
func (bot *Robot) MemoRead(key string) ([]byte, bool) {
	return bot.memo.Read(key)
}

// memorizer update
func (bot *Robot) MemoUpdate(key string, value []byte) {
	bot.memo.Update(key, value)
}

// memorizer delete
func (bot *Robot) MemoDel(key string) {
	bot.memo.Delete(key)
}

// initialize ...
func (bot *Robot) initialize() {

	// 机器人名称
	bot.name = DefaultRobotName
	if os.Getenv(`ROBOT_NAME`) != `` {
		bot.name = os.Getenv(`ROBOT_NAME`)
	}

	// 指定消息提供者，如果配置文件没有指定，则默认使用 cli
	provName := DefaultRobotProvider

	if os.Getenv(`ROBOT_PROVIDER`) != `` {
		provName = os.Getenv(`ROBOT_PROVIDER`)
	}

	prov, err := DetectProv(provName)

	if err != nil {
		panic(`Detect provider error: ` + err.Error())
	}

	bot.inMessage = prov.Incoming()
	bot.prov = prov

	// 指定储存器
	memoName := DefaultRobotMemorizer

	if os.Getenv(`ROBOT_MEMORIZER`) != `` {
		memoName = os.Getenv(`ROBOT_MEMORIZER`)
	}

	memo, err := DetectMemo(memoName)

	if err != nil {
		panic(`Detect memorizer error: ` + err.Error())
	}

	bot.memo = memo

	if bot.memo.Error() != nil {
		bot.err <- bot.memo.Error()
	}

	bot.es.init()

	bot.es.merge("custom", usrEvent)
}

