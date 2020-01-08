package dingtalk

import (
	"github.com/ghaoo/rboot"
	"github.com/hugozhu/godingtalk"
	"os"
)

type dingtalk struct {
	in     chan *rboot.Message
	out    chan *rboot.Message
	client *godingtalk.DingTalkClient
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
	}

	client := godingtalk.NewDingTalkClient(os.Getenv("DING_CORP_ID"), os.Getenv("DING_CORP_SECRET"))
	ding.client = client

	return ding
}

func init() {
	rboot.RegisterAdapter("dingtalk", newDingTalk)
}
