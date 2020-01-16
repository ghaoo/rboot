package rboot

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

const (
	// rbootLogo rboot logo
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

	version = "1.1.0"
)

// Robot 是 rboot 的一个实例，它包含了聊天转接器，规则处理器，缓存器，路由适配器和消息的进出通道
type Robot struct {
	Router     *Router
	Brain      Brain
	adapter    Adapter
	rule       Rule
	hooks      Hooks
	inputChan  chan *Message
	outputChan chan *Message

	debug      bool
	signalChan chan os.Signal
}

// New 获取一个Robot实例，
func New() *Robot {
	bot := &Robot{
		inputChan:  make(chan *Message),
		outputChan: make(chan *Message),
		signalChan: make(chan os.Signal),
		rule:       new(Regex),
	}

	bot.Router = newRouter()

	return bot
}

var processOnce sync.Once

// process 消息处理函数
func process(bot *Robot) {
	processOnce.Do(func() {

		// 监听传入消息
		for in := range bot.inputChan {

			go func(bot *Robot, msg *Message) {
				defer func() {
					if r := recover(); r != nil {
						logrus.Errorf("panic recovered when parsing message: %#v. \nPanic: %v", fmt.Sprintf(""), r)
					}
				}()

				// 匹配消息
				if script, rule, args, ok := bot.matchScript(strings.TrimSpace(msg.String())); ok {

					if bot.debug {
						logrus.Debugf("\nIncoming: \n- 类型: %s\n- 发送人: %s\n- 接收人: %v\n- 内容: %s\n- 脚本: %s\n- 规则: %s\n- 参数: %v\n\n",
							msg.Header.Get("MsgType"),
							msg.From,
							msg.To,
							msg.String(),
							script,
							rule,
							args[1:])
					}

					// 获取脚本执行函数
					action, err := DirectiveScript(script)
					if err != nil {
						logrus.Error(err)
					}

					msg.Header.Set("rule", rule)
					msg.Header["args"] = args

					// 执行脚本, 附带ctx, 并获取输出
					response := action(bot, msg)

					for _, resp := range response {
						// 将消息发送到 outputChan
						// 指定输出消息的接收者
						resp.To = msg.From

						if bot.debug {
							logrus.Debugf("\nOutgoing: \n- 类型: %s \n- 接收人: %v\n- 抄送: %v\n- 发送人: %v\n- 内容: %s\n\n",
								resp.Header.Get("MsgType"),
								resp.To,
								resp.Cc(),
								resp.From,
								resp)
						}

						// send ...
						bot.outputChan <- resp

						// 如果存在抄送人，将消息抄送给对方
						if len(resp.Cc()) > 0 {
							for _, cc := range resp.Cc() {
								resp.To = cc
								bot.outputChan <- resp
							}
						}
					}

				}

			}(bot, in)
		}
	})
}

// Go 皮皮虾，我们走~~~~~~~~~
func (bot *Robot) Go() {

	logrus.Infof("Rboot Version %s", version)

	// 初始化
	bot.initialize()

	logrus.Info("皮皮虾，我们走~~~~~~~")

	// hook
	bot.hooks.Fire(bot)

	// 开启web服务
	go bot.Router.run()

	// 消息处理
	go process(bot)

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

// Stop 皮皮虾，快停下~~~~~~~~~
func (bot *Robot) Stop() {

	runtime.SetFinalizer(bot, nil)

	logrus.Info("皮皮虾，快停下~~~~~~~~")

	os.Exit(0)
}

// Incoming 获取传入消息通道
func (bot *Robot) Incoming() chan *Message {
	return bot.inputChan
}

// Outgoing 发送消息
func (bot *Robot) Outgoing(msg *Message) {
	bot.outputChan <- msg
}

// SendText 发送文本消息
func (bot *Robot) SendText(text string, to string) {
	msg := NewMessage(text)
	msg.To = to

	bot.outputChan <- msg

}

// AddHook 添加hook
func (bot *Robot) AddHook(hook Hook) {
	bot.hooks.Add(hook)
}

// SetBrain 设置储存器
func (bot *Robot) SetBrain(brain Brain) {
	bot.Brain = brain
}

// matchScript 匹配消息内容，获取相应的脚本名称(script), 对应规则名称(matchRule), 提取的匹配内容(matchArgs)
// 当消息不匹配时，matched 返回false
func (bot *Robot) matchScript(msg string) (script, matchRule string, matchArgs []string, matched bool) {

	for script, rule := range ruleset {
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
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	bot.debug = debug

	// 指定消息提供者，如果配置文件没有指定，则默认使用 cli
	adpName := os.Getenv(`RBOOT_ADAPTER`)
	// 默认使用 cli
	if adpName == "" {
		logrus.Warn("未指定 adapter，默认使用 cli")
		adpName = "cli"
	}
	logrus.Info("已连接 ", adpName)
	adp, err := DetectAdapter(adpName)

	if err != nil {
		panic(`Detect adapter error: ` + err.Error())
	}

	// 获取转接器实例
	adapter := adp(bot)

	// 建立消息通道连接
	bot.inputChan = adapter.Incoming()
	bot.outputChan = adapter.Outgoing()

	// 储存器
	brainName := os.Getenv(`RBOOT_BRAIN`)
	// 默认使用 memory
	if brainName == "" {
		logrus.Warn("未指定 brain，默认使用 memory")
		brainName = "memory"
	}

	brain, err := DetectBrain(brainName)

	if err != nil {
		panic(`Detect brain error: ` + err.Error())
	}

	bot.Brain = brain()
}

func init() {
	_, _ = color.New(color.FgGreen).Fprintln(os.Stdout, rbootLogo)

	// 加载配置
	_ = LoadEnv()
}
