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

// Robot 封装了一个机器人运行的所有必要状态
type Robot struct {
	name     string// 机器人名称
	es       *eventStream// 事件处理器
	prov Provider// 脚本消息提供器
	conf     Config // 机器人配置信息
	Matcher string // 脚本匹配信息

	sync.Mutex
	signalChan chan os.Signal
}

// 新建Rboot机器人
func New(config ...string) *Robot {

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

// 设置名称
func (bot *Robot) SetName(name string) {
	bot.name = name
}

// 设置消息提供者
func (bot *Robot) SetProvider(provider Provider) {
	bot.prov = provider
}

// 配置
func (bot *Robot) Conf() Config {
	return bot.conf
}

// 皮皮虾，我们走~~~~~~~~~
func (bot *Robot) Go() {
	// 初始化基础信息
	bot.initialize()

	// 开启消息提供者
	go bot.prov.Run()

	// 运行事件处理器
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

// 皮皮虾，快停下......
func (bot *Robot) Stop() error {

	log.Printf("stopping %s provider", bot.prov.Name())
	if err := bot.prov.Close(); err != nil {
		return err
	}

	log.Printf("stopping %s", DefaultRobotName)
	return nil
}

func (bot *Robot) Name() string {
	return bot.name
}

// 处理消息
func (bot *Robot) handleMessage(msg Message) error {
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

// 处理 Reader 类型消息
func (bot *Robot) handleWithReader(in io.Reader) error {
	msg, err := ReadMessage(in)

	if err != nil {
		return err
	}

	return bot.handleMessage(msg)
}

// 匹配脚本规则集
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

// 匹配消息
func (bot *Robot) match(pattern, msg string) bool {

	reg := regexp.MustCompile(pattern)

	if reg.MatchString(msg) {
		return true
	}

	return false
}

// 发送消息
func (bot *Robot) Send(strs ...string) error {
	return bot.prov.Send(strs...)
}

// 回复消息
func (bot *Robot) Reply(strs ...string) error {
	return bot.prov.Reply(strs...)
}

// 初始化
func (bot *Robot) initialize() {

	// 设置名称
	if bot.conf.Name == `` {
		bot.name = DefaultRobotName
	} else {
		bot.name = bot.conf.Name
	}

	// 初始化事件处理器
	bot.es.init()

	// 加载自定义事件处理器
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

	bot.prov = prov(bot)
}
