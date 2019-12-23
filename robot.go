package rboot

import (
	"context"
	"github.com/ghaoo/rboot/tools/env"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"sync"
)

var AppName string

const Version = "3.1.0"

type Robot struct {
	adapter  Adapter
	rule     Rule
	contacts []User
	hooks    []Hook
	MatchRule string
	MatchSub []string

	inputChan  chan Message
	outputChan chan Message

	sync.RWMutex

	conf Config
}

// New 获取Robot实例
func New() *Robot {

	bot := &Robot{
		inputChan:  make(chan Message),
		outputChan: make(chan Message),
		conf:       newConfig(),
		rule:       new(Regex),
	}

	return bot
}

var processOnce sync.Once

// process 消息处理函数
func process(ctx context.Context, bot *Robot) {
	processOnce.Do(func() {

		// 监听传入消息
		for in := range bot.inputChan {
			go func(bot Robot, msg Message) {
				defer func() {
					if r := recover(); r != nil {
						logrus.Errorf("panic recovered when parsing message: %#v. \nPanic: %v", msg, r)
					}
				}()

				// 将传入消息拷贝到 ctx
				ctx = context.WithValue(ctx, "input", msg)

				// 匹配消息
				if script, mr, ms, ok := bot.matchScript(msg.Content); ok {

					// 匹配的脚本对应规则
					bot.MatchRule = mr

					// 消息匹配集合
					bot.MatchSub = ms

					// 获取脚本执行函数
					action, err := DirectiveScript(script)

					if err != nil {
						logrus.Error(err)
					}

					// 执行脚本并获取输出，附带 ctx
					responses := action(ctx, &bot)

					// 将消息发送到 outputChan
					if len(responses) > 0 {
						for _, resp := range responses {
							// 指定输出消息的接收者和发送者
							resp.From = msg.To
							resp.To = msg.From

							// send ...
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
	// 设置Robot名称
	AppName = bot.conf.Name

	// 初始化
	bot.initialize()

	logrus.Debug(`皮皮虾，我们走~~~~~~~`)

	// 上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = context.WithValue(ctx, "appname", AppName)

	// 消息处理
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

// SyncUsers 同步用户
func (bot *Robot) SyncUsers(user []User) {
	bot.Lock()
	if len(user) > 0 {
		bot.contacts = user
	}
	bot.Unlock()
}

// 消息入站
func (bot *Robot) Incoming(msg Message) {
	bot.Lock()
	bot.inputChan <- msg
	bot.Unlock()
}

func (bot *Robot) Outgoing() chan Message {
	bot.Lock()
	defer bot.Unlock()

	return bot.outputChan
}

// Send 发送消息
func (bot *Robot) Send(msg Message) {
	bot.outputChan <- msg
}

// SendText 发送文本消息
func (bot *Robot) SendText(text string, to ...User) {

	if len(to) > 0 {
		for _, user := range to {
			msg := Message{
				To:      user,
				Content: text,
			}
			bot.outputChan <- msg
		}
	} else {
		bot.outputChan <- Message{Content: text}
	}

}

// 直接执行脚本命令
/*func (bot *Robot) RunScript(script, cmd string) error {

	return nil
}*/

// MatchScript 匹配消息内容，获取相应的脚本名称(script), 对应规则名称(matchRule), 提取的匹配内容(match)
// 当消息不匹配时，matched 返回false
func (bot *Robot) matchScript(msg string) (script, matchRule string, match []string, matched bool) {

	for script, rule := range rulesets {
		for m, r := range rule {
			if match, ok := bot.rule.Match(r, msg); ok {
				return script, m, match, true
			}
		}
	}

	return ``, ``, nil, false
}

// initialize 初始化 Robot
func (bot *Robot) initialize() {

	// 指定消息提供者，如果配置文件没有指定，则默认使用 cli
	adp, err := DetectAdapter(bot.conf.Adapter)

	if err != nil {
		panic(`Detect adapter error: ` + err.Error())
	}

	// 获取适配器实例
	adapter := adp(bot)

	// 建立消息通道连接
	bot.inputChan = adapter.Incoming()
	bot.outputChan = adapter.Outgoing()
}

func init() {

	// 加载配置
	err := env.Load()

	if err != nil {
		panic(`Load env config error: ` + err.Error())
	}
}
