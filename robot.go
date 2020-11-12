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

	version = "1.2.7"
)

var defaultCachePath = ".data"

// Robot 是 rboot 的一个实例，它包含了聊天转接器，规则处理器，缓存器，路由适配器和消息的进出通道
type Robot struct {
	// 路由，支持脚本自定义添加新路由
	Router *Router
	// 缓存
	Brain Brain
	// 钩子
	Hooks MsgHooks
	// 缓存文件夹
	CachePath string
	// 调试信息
	Debug bool

	// 终端
	adapter Adapter
	// 传入消息
	inputChan chan *Message
	// 传出消息
	outputChan chan *Message

	// 规则匹配器
	rule Rule
	// 操作系统信号
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

	bot.CachePath = defaultCachePath
	if len(os.Getenv("CACHE_PATH")) > 0 {
		bot.CachePath = os.Getenv("CACHE_PATH")
	}

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	bot.Debug = debug

	bot.Router = newRouter()

	// 初始化
	bot.initialize()

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
						logrus.WithFields(logrus.Fields{
							"mod": "rboot",
							"msg": msg,
						}).Errorf("panic recovered when parsing message: %#v. \nPanic: %v", msg, r)
					}
				}()

				// 处理消息前的Hook
				bot.fireHooks(HOOK_BEFORE_INCOMING, msg)

				// 匹配消息
				if plug, rule, args, ok := bot.matchPlugin(strings.TrimSpace(msg.String())); ok {

					if bot.Debug {
						logrus.Debugf("- 脚本: %s\n- 规则: %s\n- 参数: %v\n\n",
							plug,
							rule,
							args[1:])
					}

					// 获取插件执行函数
					action, err := DirectivePlugin(plug)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"mod":     "rboot",
							"plug":    plug,
							"ruleset": rule,
							"msg":     msg,
						}).WithError(err).Error("listen: directive plugin error")
					}

					msg.Header.Set("rule", rule)
					msg.Header["args"] = args

					// 执行脚本, 附带ctx, 并获取输出
					response := action(bot, msg)

					for _, resp := range response {
						// 将消息发送到 outputChan
						// 指定输出消息的接收者
						resp.To = msg.From

						if msg.KeepHeader {
							for hn, hv := range msg.Header {
								if len(resp.Header[hn]) <= 0 {
									resp.Header[hn] = hv
								}
							}
						}

						if bot.Debug {
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
						// 处理消息后的Hook
						bot.fireHooks(HOOK_AFTER_OUTGOING, msg)
					}
				}

			}(bot, in)
		}
	})
}

// Go 皮皮虾，我们走~~~~~~~~~
func (bot *Robot) Go() {

	fmt.Println("Rboot Version ", version)

	fmt.Println("皮皮虾，我们走~~~~~~~")

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

	fmt.Println("皮皮虾，快停下~~~~~~~~")

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

// SetBrain 设置储存器
func (bot *Robot) SetBrain(brain Brain) {
	bot.Brain = brain
}

// fireHooks 执行钩子
func (bot *Robot) fireHooks(typ int, msg *Message) {
	err := bot.Hooks.Fire(typ, bot, msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"mod":  "rboot",
			"type": hookType[typ],
			"msg":  msg,
		}).WithError(err).Error("fireHooks: fire hooks failed")
	}
}

// matchScript 匹配消息内容，获取相应的脚本名称(script), 对应规则名称(matchRule), 提取的匹配内容(matchArgs)
// 当消息不匹配时，matched 返回false
func (bot *Robot) matchPlugin(msg string) (script, matchRule string, matchArgs []string, matched bool) {

	for script, rules := range ruleset {
		for m, r := range rules {
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
	adpName := os.Getenv(`ROBOT_ADAPTER`)
	// 默认使用 cli
	if adpName == "" {
		fmt.Println("未指定 adapter，默认使用 cli")
		adpName = "cli"
	}
	fmt.Println("已连接适配器 ", adpName)
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
	brainName := os.Getenv(`ROBOT_BRAIN`)
	// 默认使用 memory
	if brainName == "" {
		brainName = "bolt"
	}

	brain, err := DetectBrain(brainName)

	if err != nil {
		panic(`Detect brain error: ` + err.Error())
	}

	bot.Brain = brain()

	// 监听 http 入站消息的 ResultFul API
	bot.Router.HandleFunc("/incoming", bot.listenIncoming).Methods("POST")
}

func init() {
	_, _ = color.New(color.FgGreen).Fprintln(os.Stdout, rbootLogo)

	// 加载配置
	_ = LoadEnv()
}
