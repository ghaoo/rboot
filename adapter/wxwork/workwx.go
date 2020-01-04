package adapter

import (
	"github.com/ghaoo/rboot"
	"github.com/ghaoo/wxwork"
)

type wework struct {
	in  chan rboot.Message
	out chan rboot.Message
}

func newWework(bot *rboot.Robot) rboot.Adapter {
	w := &wework{
		in:  make(chan rboot.Message),
		out: make(chan rboot.Message),
	}

	agent := agent()
	bot.Router.HandleFunc("/wxwork", agent.CallbackVerify)

	return w
}

func agent() *wxwork.Agent {
	a := wxwork.NewAgent("ww7a4b068e86da007a", "Qvy3I_NYQsRDXp8btU9ips4BuflVmZUtMUiDPH4I2Rg", 1000002)
	a.SetCallback("6OxE5PFldOqKnilqC6CWAH", "Ke6rh3wv1KJvJOJcaeOVL41Y54AN2KwiIPHq3DMxNDo")

	return a
}

func (w *wework) Name() string {
	return "wework"
}

func (w *wework) Incoming() chan rboot.Message {
	return w.in
}

func (w *wework) Outgoing() chan rboot.Message {
	return w.out
}

// 监听 rboot Outgoing
func (w *wework) listenOutgoing() {}

func init() {
	rboot.RegisterAdapter("wxwork", newWework)
}
