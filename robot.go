package rboot

import (
	"context"
	"github.com/ghaoo/rboot/tools/env"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
)

var AppName string

const Version = "3.1.0"

type Robot struct {
	adapter    Adapter
	brain      Brain
	rule       Rule
	history    Histories
	contacts   []User
	inputChan  chan Message
	outputChan chan Message

	MatchRule string
	MatchSub  []string

	signalChan chan os.Signal
	sync.RWMutex
}

// New 获取Robot实例
func New() *Robot {

	bot := &Robot{
		history: Histories{},
		inputChan:  make(chan Message),
		outputChan: make(chan Message),
		signalChan: make(chan os.Signal),
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

					// 执行脚本, 附带ctx, 并获取输出
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

					// 将用户操作写入历史
					rh, _ := strconv.ParseBool(os.Getenv(`RECORD_HISTORY`))
					if rh {
						bot.history.Push(msg, responses)
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
	AppName = os.Getenv(`RBOOT_NAME`)

	// 初始化
	bot.initialize()

	logrus.Debug(`皮皮虾，我们走~~~~~~~`)

	// 上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = context.WithValue(ctx, "appname", AppName)

	// 消息处理
	go process(ctx, bot)

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

	runtime.SetFinalizer(bot, nil)

	logrus.Debug(`皮皮虾，快停下~~~`)

	os.Exit(0)

	return nil
}

// SetRule 设置消息处理器
func (bot *Robot) SetRule(rule Rule) {
	bot.Lock()
	bot.rule = rule
	bot.Unlock()
}

// SyncUsers 同步用户
func (bot *Robot) SyncUsers(user []User) {
	bot.Lock()
	if len(user) > 0 {
		bot.contacts = user
	}
	bot.Unlock()
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

// 设置储存器
func (bot *Robot) SetBrain(brain Brain) {
	bot.Lock()
	bot.brain = brain
	bot.Unlock()
}

// Brain set ...
func (bot *Robot) BrainSet(key string, value []byte) error {
	return bot.brain.Set(key, value)
}

// Brain get ...
func (bot *Robot) BrainGet(key string) []byte {
	return bot.brain.Get(key)
}

// Brain get ...
func (bot *Robot) BrainRemove(key string) error {
	return bot.brain.Remove(key)
}

// 上一条信息（历史记录）
func (bot *Robot) PrevHistory(uid string) *History {
	bot.Lock()
	defer bot.Unlock()
	return bot.history.Prev(uid)
}

// 前几条信息（历史记录）
func (bot *Robot) PrevHistoryN(uid string, n int) []*History {
	bot.Lock()
	defer bot.Unlock()
	return bot.history.PrevN(uid, n)
}

// 清空历史消息
func (bot *Robot) ClearHistory(uid string) {
	bot.Lock()
	defer bot.Unlock()
	bot.history.Clear(uid)
}

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
	adp_name := os.Getenv(`RBOOT_ADAPTER`)
	// 默认使用 cli
	if adp_name == "" {
		adp_name = "cli"
	}
	adp, err := DetectAdapter(adp_name)

	if err != nil {
		panic(`Detect adapter error: ` + err.Error())
	}

	// 获取适配器实例
	adapter := adp(bot)

	// 建立消息通道连接
	bot.inputChan = adapter.Incoming()
	bot.outputChan = adapter.Outgoing()

	// 储存器
	brain_name := os.Getenv(`RBOOT_BRAIN`)
	// 默认使用 memory
	if brain_name == "" {
		brain_name = "memory"
	}
	brain, err := DetectBrain(brain_name)

	if err != nil {
		panic(`Detect brain error: ` + err.Error())
	}

	bot.brain = brain()
}

func init() {

	// 加载配置
	err := env.Load()

	if err != nil {
		logrus.Error(`Load env config error: `, err.Error())
	}
}
