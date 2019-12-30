package rboot

import (
	"context"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/ghaoo/rboot/tools/env"
	"github.com/sirupsen/logrus"
)

var AppName string

const (
	rbootLogo = `
===================================================================
*   ________  ____  ____  ____  ______   ________  ____  ______   *
*   ___/ __ \/ __ )/ __ \/ __ \/_  __/   ___/ __ )/ __ \/_  __/   *
*   __/ /_/ / __  / / / / / / / / /      __/ __  / / / / / /      *
*   _/ _  _/ /_/ / /_/ / /_/ / / /       _/ /_/ / /_/ / / /       *
*   /_/ |_/_____/\____/\____/ /_/        /_____/\____/ /_/        *
*                                                                 *
*                      Powerful and Happy                         *
===================================================================
`

	Version = "1.1.0"
)

type Robot struct {
	adapter    Adapter
	brain      Brain
	rule       Rule
	contacts   []User
	inputChan  chan Message
	outputChan chan Message

	Router  *Router
	Ruleset string
	Args    []string

	signalChan chan os.Signal
	mu         sync.RWMutex
}

// New 获取Robot实例
func New() *Robot {

	bot := &Robot{
		inputChan:  make(chan Message),
		outputChan: make(chan Message),
		signalChan: make(chan os.Signal),
		rule:       new(Regex),
	}

	bot.Router = newRouter()

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

				msg.Content = strings.TrimSpace(msg.Content)

				// 将传入消息拷贝到 ctx
				ctx = context.WithValue(ctx, "input", msg)

				// 匹配消息
				if script, mr, ms, ok := bot.matchScript(msg.Content); ok {

					// 匹配的脚本对应规则
					bot.Ruleset = mr

					// 消息匹配参数
					bot.Args = ms

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

							if resp.Channel == "" {
								resp.Channel = msg.Channel
							}

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

	logrus.Infof("Rboot Version %s", Version)
	// 设置Robot名称
	AppName = os.Getenv(`RBOOT_NAME`)

	// 初始化
	bot.initialize()

	logrus.Info(`皮皮虾，我们走~~~~~~~`)

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
func (bot *Robot) Stop() {

	runtime.SetFinalizer(bot, nil)

	logrus.Info(`皮皮虾，快停下~~~`)

	os.Exit(0)
}

// SyncUsers 同步用户
func (bot *Robot) SyncUsers(user []User) {
	bot.mu.Lock()
	if len(user) > 0 {
		bot.contacts = user
	}
	bot.mu.Unlock()
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
	bot.mu.Lock()
	bot.brain = brain
	bot.mu.Unlock()
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
	adpName := os.Getenv(`RBOOT_ADAPTER`)
	// 默认使用 cli
	if adpName == "" {
		adpName = "cli"
	}
	adp, err := DetectAdapter(adpName)

	if err != nil {
		panic(`Detect adapter error: ` + err.Error())
	}

	// 获取适配器实例
	adapter := adp(bot)

	// 建立消息通道连接
	bot.inputChan = adapter.Incoming()
	bot.outputChan = adapter.Outgoing()

	// 储存器
	brainName := os.Getenv(`RBOOT_BRAIN`)
	// 默认使用 memory
	if brainName == "" {
		brainName = "memory"
	}
	brain, err := DetectBrain(brainName)

	if err != nil {
		panic(`Detect brain error: ` + err.Error())
	}

	bot.brain = brain()

	// 开启web服务
	go bot.Router.run()
}

func init() {
	color.New(color.FgGreen).Fprintln(os.Stdout, rbootLogo)

	// 加载配置
	err := env.Load()

	if err != nil {
		logrus.Error(`Load env config error: `, err.Error())
	}
}
