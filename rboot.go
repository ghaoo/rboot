package rboot

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
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

	providerIn  chan Message
	providerOut chan Message

	signalChan chan os.Signal
	sync.Mutex
}

func New() *Robot {

	bot := &Robot{
		es:          newStream(),
		providerIn:  make(chan Message),
		providerOut: make(chan Message),
		signalChan:  make(chan os.Signal, 1),
	}

	return bot
}

func (bot *Robot) SetName(name string) {
	bot.name = name
}

func (bot *Robot) Send(msg Message) {
	bot.providerOut <- msg
}

var processOnce sync.Once

func (bot *Robot) process() {
	processOnce.Do(func() {

		for _, script := range scripts {
			script.Hook(*bot)
		}

		for in := range bot.providerIn {
			go func(bot Robot, msg Message) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("panic recovered when parsing message: %#v. Panic: %v", msg, r)
					}
				}()

				for _, script := range scripts {
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

// 皮皮虾，快停下~~~~~~~~~
func (bot *Robot) Stop() error {

	log.Printf("stopping %s", DefaultRobotName)
	return nil
}

// robot name
func (bot *Robot) Name() string {
	return bot.name
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

	bot.providerIn = prov.Incoming()
	bot.providerOut = prov.Outgoing()

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
		log.Print(bot.memo.Error())
	}

	bot.es.init()

	bot.es.merge("custom", usrEvent)
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
