package dingtalk

import (
	"github.com/ghaoo/rboot"
	"os"
)

type dingtalk struct {
	secret string

	in  chan *rboot.Message
	out chan *rboot.Message
	bot *rboot.Robot
}

func (ding *dingtalk) Name() string {
	return "dingtalk"
}

func (ding *dingtalk) Incoming() chan *rboot.Message {
	return ding.in
}

func (ding *dingtalk) Outgoing() chan *rboot.Message {
	return ding.out
}

func newDingTalk(bot *rboot.Robot) rboot.Adapter {
	ding := &dingtalk{
		in:  make(chan *rboot.Message),
		out: make(chan *rboot.Message),
		bot: bot,
	}

	secret := os.Getenv("DING_ROBOT_SECRET")
	if secret == "" {
		panic("DING_ROBOT_SECRET 未设置!!")
	}

	ding.secret = secret

	bot.Router.HandleFunc("/ding", ding.listenIncoming).Methods("POST")

	// 当需要定时推送消息或警报消息时将监听 outgoing 打开
	// 创建自定义机器人，并设置 DING_ROBOT_HOOK_ACCESS_TOKEN
	// 和 DING_ROBOT_HOOK_SECRET，所有消息都由自定义机器人回复
	go ding.listenOutgoing()

	return ding
}

func init() {
	rboot.RegisterAdapter("dingtalk", newDingTalk)
}
